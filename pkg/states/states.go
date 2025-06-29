package states

import (
	"encoding/json"
	"fmt"
)

type State struct {
	Sigla    string
	LongName string
	Region   string
}

func (s *State) UnmarshalJSON(data []byte) error {
	var sigla string
	if err := json.Unmarshal(data, &sigla); err != nil {
		return err
	}

	state, exists := States[sigla]
	if !exists {
		return fmt.Errorf("state %s does not exist", sigla)
	}

	*s = state
	return nil
}

func (s *State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Sigla)
}

var (
	AC = State{Sigla: "AC", LongName: "Acre", Region: "Norte"}
	AL = State{Sigla: "AL", LongName: "Alagoas", Region: "Nordeste"}
	AP = State{Sigla: "AP", LongName: "Amapá", Region: "Norte"}
	AM = State{Sigla: "AM", LongName: "Amazonas", Region: "Norte"}
	BA = State{Sigla: "BA", LongName: "Bahia", Region: "Nordeste"}
	CE = State{Sigla: "CE", LongName: "Ceará", Region: "Nordeste"}
	ES = State{Sigla: "ES", LongName: "Espírito Santo", Region: "Sudeste"}
	GO = State{Sigla: "GO", LongName: "Goiás", Region: "Centro-Oeste"}
	MA = State{Sigla: "MA", LongName: "Maranhão", Region: "Nordeste"}
	MT = State{Sigla: "MT", LongName: "Mato Grosso", Region: "Centro-Oeste"}
	MS = State{Sigla: "MS", LongName: "Mato Grosso do Sul", Region: "Centro-Oeste"}
	MG = State{Sigla: "MG", LongName: "Minas Gerais", Region: "Sudeste"}
	PA = State{Sigla: "PA", LongName: "Pará", Region: "Norte"}
	PB = State{Sigla: "PB", LongName: "Paraíba", Region: "Nordeste"}
	PR = State{Sigla: "PR", LongName: "Paraná", Region: "Sul"}
	PE = State{Sigla: "PE", LongName: "Pernambuco", Region: "Nordeste"}
	PI = State{Sigla: "PI", LongName: "Piauí", Region: "Nordeste"}
	RJ = State{Sigla: "RJ", LongName: "Rio de Janeiro", Region: "Sudeste"}
	RN = State{Sigla: "RN", LongName: "Rio Grande do Norte", Region: "Nordeste"}
	RS = State{Sigla: "RS", LongName: "Rio Grande do Sul", Region: "Sul"}
	RO = State{Sigla: "RO", LongName: "Rondônia", Region: "Norte"}
	RR = State{Sigla: "RR", LongName: "Roraima", Region: "Norte"}
	SC = State{Sigla: "SC", LongName: "Santa Catarina", Region: "Sul"}
	SP = State{Sigla: "SP", LongName: "São Paulo", Region: "Sudeste"}
	SE = State{Sigla: "SE", LongName: "Sergipe", Region: "Nordeste"}
	TO = State{Sigla: "TO", LongName: "Tocantins", Region: "Norte"}
	DF = State{Sigla: "DF", LongName: "Distrito Federal", Region: "Centro-Oeste"}
)

var States = map[string]State{
	"AC": AC,
	"AL": AL,
	"AP": AP,
	"AM": AM,
	"BA": BA,
	"CE": CE,
	"DF": DF,
	"ES": ES,
	"GO": GO,
	"MA": MA,
	"MT": MT,
	"MS": MS,
	"MG": MG,
	"PA": PA,
	"PB": PB,
	"PR": PR,
	"PE": PE,
	"PI": PI,
	"RJ": RJ,
	"RN": RN,
	"RS": RS,
	"RO": RO,
	"RR": RR,
	"SC": SC,
	"SP": SP,
	"SE": SE,
	"TO": TO,
}

type Region struct {
	Name   string
	States []State
}

var (
	Sul         = Region{Name: "Sul", States: []State{PR, RS, SC}}
	CentroOeste = Region{Name: "Centro-Oeste", States: []State{DF, GO, MT, MS}}
	Nordeste    = Region{Name: "Nordeste", States: []State{AL, BA, CE, MA, PB, PE, PI, RN, SE}}
	Norte       = Region{Name: "Norte", States: []State{AC, AP, AM, PA, RO, RR, TO}}
	Sudeste     = Region{Name: "Sudeste", States: []State{ES, MG, RJ, SP}}
)

var Regions = map[string]Region{
	"Sul":          Sul,
	"Centro-Oeste": CentroOeste,
	"Nordeste":     Nordeste,
	"Norte":        Norte,
	"Sudeste":      Sudeste,
}

func (r *Region) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	region, exists := Regions[name]
	if !exists {
		return fmt.Errorf("state %s does not exist", name)
	}

	*r = region
	return nil
}

func (r *Region) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Name)
}
