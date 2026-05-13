package plategenerator

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"encomers/license/internal/domain/services/generator"
	"encomers/license/internal/domain/services/repository"
	models "encomers/license/internal/domain/valueObjects"
)

type PlateGenerator struct {
	numberGen   generator.INumberGenerator
	letterGen   generator.ILetterGenerator
	regionGen   generator.IRegionGenerator
	numbers     []models.NumberPart
	letters     []models.LetterPart
	repository  repository.IRepository
	numberIndex int
	letterIndex int
	used        int
	total       int
	pairs       map[int]map[int]struct{}
	overflow    bool
	mu          sync.Mutex
}

func New(numberGen generator.INumberGenerator, letterGen generator.ILetterGenerator, regionGen generator.IRegionGenerator, repository repository.IRepository) (*PlateGenerator, error) {
	if numberGen == nil || letterGen == nil || regionGen == nil || repository == nil {
		return nil, fmt.Errorf("plate generator parameters are invalid")
	}

	letters := letterGen.Size()
	letterCombinations := 1
	for range models.LETTERS_COUNT {
		letterCombinations *= letters
	}

	result := &PlateGenerator{
		numberGen:   numberGen,
		letterGen:   letterGen,
		regionGen:   regionGen,
		repository:  repository,
		pairs:       make(map[int]map[int]struct{}),
		letters:     make([]models.LetterPart, 0, letterCombinations),
		numbers:     make([]models.NumberPart, 0, models.MAX_NUMBER_PART_VALUE+1),
		numberIndex: 0,
		letterIndex: 0,
		used:        0,
		total:       0,
		overflow:    false,
	}

	err := result.init()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PlateGenerator) init() error {
	p.numbers = p.numbers[:0]
	p.letters = p.letters[:0]

	for !p.numberGen.IsOverflow() {
		next := p.numberGen.GetNext()
		if p.numberGen.IsOverflow() {
			break
		}
		p.numbers = append(p.numbers, next)
	}
	for !p.letterGen.IsOverflow() {
		next := p.letterGen.GetNext()
		if p.letterGen.IsOverflow() {
			break
		}
		p.letters = append(p.letters, next)
	}

	if len(p.numbers) == 0 || len(p.letters) == 0 {
		p.overflow = true
		return nil
	}
	return p.rebuildPairs()
}

func (p *PlateGenerator) rebuildPairs() error {
	p.pairs = make(map[int]map[int]struct{}, len(p.letters))

	p.total = len(p.numbers) * len(p.letters)
	region := p.regionGen.Get()

	for i, letter := range p.letters {
		row := make(map[int]struct{})
		p.pairs[i] = row

		for j, number := range p.numbers {
			plate := models.Plate{
				Numbers: number,
				Letters: letter,
				Region:  region,
			}

			ok, err := p.repository.CheckPlate(plate)
			if err != nil {
				return err
			}

			if ok {
				row[j] = struct{}{}
				p.used++
			}
		}
	}

	p.overflow = p.used >= p.total

	return nil
}

func (p *PlateGenerator) IsOverflow() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.overflow
}

func (p *PlateGenerator) inc() {
	if p.overflow {
		return
	} else if _, ok := p.pairs[p.letterIndex][p.numberIndex]; !ok {
		return
	} else {
		startLetterIndex := p.letterIndex
		startNumberIndex := p.numberIndex
		for !p.overflow {
			p.numberIndex++
			if p.numberIndex == len(p.numbers) {
				p.numberIndex = 0
				p.letterIndex++
				if p.letterIndex == len(p.letters) {
					p.letterIndex = 0

				}
			}

			if _, ok := p.pairs[p.letterIndex][p.numberIndex]; !ok {
				return
			}

			if p.letterIndex == startLetterIndex && p.numberIndex == startNumberIndex {
				p.overflow = true
				return
			}
		}
	}
}

func (p *PlateGenerator) getCurrent() models.Plate {
	return models.Plate{
		Numbers: p.numbers[p.numberIndex],
		Letters: p.letters[p.letterIndex],
		Region:  p.regionGen.Get(),
	}
}

func (p *PlateGenerator) GetNext() (models.Plate, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.inc()

	if p.overflow {
		return models.Plate{}, fmt.Errorf("no more plates available")
	}

	plate := p.getCurrent()
	err := p.repository.SavePlate(plate)
	if err != nil {
		return models.Plate{}, err
	}

	p.pairs[p.letterIndex][p.numberIndex] = struct{}{}
	p.used++

	return plate, nil
}

func (p *PlateGenerator) GetRandom() (models.Plate, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.overflow {
		return models.Plate{}, fmt.Errorf("no more plates available")
	}
	for {
		p.letterIndex = rand.IntN(len(p.letters))
		p.numberIndex = rand.IntN(len(p.numbers))
		if _, ok := p.pairs[p.letterIndex][p.numberIndex]; !ok {
			break
		}
	}

	plate := p.getCurrent()
	err := p.repository.SavePlate(plate)
	if err != nil {
		return models.Plate{}, err
	}

	p.pairs[p.letterIndex][p.numberIndex] = struct{}{}
	p.used++
	return plate, nil
}
