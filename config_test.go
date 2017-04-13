// config_test
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gcfg.v1"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	err := gcfg.ReadFileInto(&cfg, "oceano2oceansites.toml")

	// test if file exist
	assert.Nil(err)
	// test Global section
	assert.Equal(cfg.Global.Author, "jgrelet IRD july 2015 CASSIOPEE cruise")
	assert.Equal(cfg.Global.Debug, false)
	assert.Equal(cfg.Global.Echo, true)
	// test Cruise section
	assert.Equal(cfg.Cruise.CycleMesure, "CASSIOPEE")
	assert.Equal(cfg.Cruise.Plateforme, "ATALANTE")
	assert.Equal(cfg.Cruise.Callsign, "FNCM")
	assert.Equal(cfg.Cruise.Institute, "IRD")
	assert.Equal(cfg.Cruise.Pi, "MARIN")
	assert.Equal(cfg.Cruise.Timezone, "GMT")
	assert.Equal(cfg.Cruise.BeginDate, "19/07/2015")
	assert.Equal(cfg.Cruise.EndDate, "23/08/2015")
	assert.Equal(cfg.Cruise.Creator, "Jacques.Grelet_at_ird.fr")
	// test ctd section
	assert.Equal(cfg.Ctd.CruisePrefix, "csp")
	assert.Equal(cfg.Ctd.StationPrefixLength, 5)
	assert.Equal(cfg.Ctd.TitleSummary, "CTD profiles processed during CASSIOPEE cruise")
	assert.Equal(cfg.Ctd.TypeInstrument, "SBE911+")
	assert.Equal(cfg.Ctd.InstrumentNumber, "09P29544-0694")
	assert.Equal(cfg.Ctd.Split, "PRES,3,DEPTH,4,ETDD,2,TEMP,5,PSAL,22,DENS,24,SVEL,26,DOX2,19,FLU2,14,TUR3,13,LGH3,15,NUMP,18,NAVG,21")
	assert.Equal(cfg.Ctd.SplitAll, "PRES,3,DEPTH,4,ETDD,2,TE01,5,TE02,6,PSA1,22,PSA2,23,DO12,19,DO22,20,DEN1,24,DEN2,25,SVEL,26,CND1,7,CND2,8,DOV1,9,DVT1,10,DOV2,11,DVT2,12,TUR3,13,FLU2,14,LGH3,15,LGHT,16,LGH4,17,NUMP,18,NAVG,21")
	// test section btl
	assert.Equal(cfg.Btl.CruisePrefix, "csp")
	assert.Equal(cfg.Btl.StationPrefixLength, 5)
	assert.Equal(cfg.Btl.TitleSummary, "Water sample during cruise with 22 levels")
	assert.Equal(cfg.Btl.TypeInstrument, "SBE32 standard 22 Niskin bottles")
	assert.Equal(cfg.Btl.InstrumentNumber, "unknown")
	assert.Equal(cfg.Btl.Comment, "CTD bottles water sampling with temperature, salinity and oxygen from primary and secondary sensors")
	assert.Equal(cfg.Btl.Split, "BOTL,1,PSA1,5,PSA2,6,DO11,7,DO21,8,PRES,14,DEPTH,15")
}
