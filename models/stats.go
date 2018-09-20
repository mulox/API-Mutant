package models

type Stats struct{
	Mutant int `json:"count_mutant_dna"`
	Human int `json:"count_human_dna"`
	Ratio float32 `json:"ratio"`
}