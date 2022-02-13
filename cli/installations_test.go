package cli

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"github.com/satisfactorymodding/ficsit-cli/cfg"
	"github.com/satisfactorymodding/ficsit-cli/tea/utils"
)

func init() {
	cfg.SetDefaults()
}

func TestInstallationsInit(t *testing.T) {
	installations, err := InitInstallations()
	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, installations)
}

func TestAddInstallation(t *testing.T) {
	ctx, err := InitCLI()
	testza.AssertNoError(t, err)

	profileName := "InstallationTest"
	profile, err := ctx.Profiles.AddProfile(profileName)
	testza.AssertNoError(t, err)
	testza.AssertNoError(t, profile.AddMod(utils.Mod{
		Name:      "Area Actions",
		ID:        "6vQ6ckVYFiidDh",
		Reference: "AreaActions"}, ">=1.6.5"))
	testza.AssertNoError(t, profile.AddMod(utils.Mod{
		Name: "ArmorModules__Modpack_All"}, ">=1.4.1"))

	installation, err := ctx.Installations.AddInstallation(ctx, "../testdata/server", profileName)
	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, installation)

	err = installation.Install(ctx)
	testza.AssertNoError(t, err)
}
