package ficsit

import (
	"context"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/MarvinJWendt/testza"
	"github.com/satisfactorymodding/ficsit-cli/cfg"
)

var client graphql.Client

func init() {
	cfg.SetDefaults()
	client = InitAPI()
}

func TestModVersions(t *testing.T) {
	response, err := ModVersions(context.Background(), client, "SmartFoundations", VersionFilter{})
	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, response)
	testza.AssertNotNil(t, response.GetMod)
	testza.AssertNotNil(t, response.GetMod.Versions)
	testza.AssertNotZero(t, len(response.GetMod.Versions))
}

func TestMods(t *testing.T) {
	response, err := Mods(context.Background(), client, ModFilter{})
	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, response)
	testza.AssertNotNil(t, response.GetMods)
	testza.AssertNotNil(t, response.GetMods.Mods)
	testza.AssertNotZero(t, response.GetMods.Count)
	testza.AssertNotZero(t, len(response.GetMods.Mods))
}
