package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/spigell/pulumi-automation-api-apps/hetzner-snaphots-manager/sdk/snapshots"
)

type configuration []Machine

type Machine struct {
    ID string
}


var (
        defaultImageName = "ubuntu-20.04"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Get default ID for ubuntu-20.04 image
        defaultImage, err := hcloud.GetImage(ctx, &hcloud.GetImageArgs{
            Name: &defaultImageName,
        })

        if err != nil {
            return fmt.Errorf("get default image: %w", err)
        }

        // Get config for stack
        cfg := config.New(ctx, "")
        var pulumiConfig configuration
        cfg.RequireObject("machines", &pulumiConfig)

        for _, machine := range pulumiConfig {
            sn, err := snapshots.GetLastSnapshot(&http.Client{}, machine.ID)

            if err != nil {
                switch {
                case errors.Is(err, snapshots.ErrSnapshotNotFound):
                    sn.Body.ID = defaultImage.Id
                default:
                    return fmt.Errorf("get uncovered error for last snapshot: %w", err)
                }
            }

            // Create a new Hetzner Cloud server
            _, err = hcloud.NewServer(ctx, machine.ID, &hcloud.ServerArgs{
                Name:       pulumi.String(machine.ID),
                ServerType: pulumi.String("cx11"),
                Image:      pulumi.String(strconv.Itoa(sn.Body.ID)),
                Location:   pulumi.String("nbg1"),
                UserData: pulumi.String(`
                #!/bin/bash
                echo "Hello from Pulumi on Hetzner Cloud!" > /root/hello.txt
                `),
            })

            if err != nil {
                return err
            }
        }

        return nil
    })
}
