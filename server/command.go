package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const cocoCommand = "coco"

func createCocoCommand(cocoCategories map[string][]string) *model.Command {
	return &model.Command{
		Trigger:          cocoCommand,
		AutoComplete:     true,
		AutoCompleteDesc: "Draw a textmoji",
		AutoCompleteHint: "[name]",
		AutocompleteData: getAutocompleteData(cocoCategories),
	}
}

func normalizeName(name string) string {
	result := strings.ToLower(name)
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, ".svg", "")
	return result
}

func nameToUrl(name string) string {
	return "/plugins/com.cocomaterial/" + normalizeName(name) + ".png"
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	action := ""
	if len(split) > 1 {
		action = split[1]
	}

	if command != "/"+cocoCommand {
		return &model.CommandResponse{}, nil
	}

	for _, cocoEntry := range p.cocoEntries {
		if action == normalizeName(cocoEntry) {
			p.API.CreatePost(&model.Post{
				Message:   fmt.Sprintf("![%s](%s =x24)", normalizeName(cocoEntry), nameToUrl(cocoEntry)) + strings.TrimPrefix(args.Command, command+" "+action),
				UserId:    args.UserId,
				ChannelId: args.ChannelId,
				ParentId:  args.ParentId,
				RootId:    args.RootId,
			})
			return &model.CommandResponse{}, nil
		}
	}
	return &model.CommandResponse{}, nil
}

func getAutocompleteData(cocoCategories map[string][]string) *model.AutocompleteData {
	coco := model.NewAutocompleteData("coco", "[name] [extra-text]", "Draw a coco material and add an extra text later")
	coco.CanAutocompleteInTheMiddle = true

	for cocoCategoryName, cocoEntries := range cocoCategories {
		subcatcommand := model.NewAutocompleteData(cocoCategoryName, "coco [extra-text]", fmt.Sprintf("Draw the coco material from the %s category", cocoCategoryName))
		subcatcommand.CanAutocompleteInTheMiddle = true
		coco.AddCommand(subcatcommand)
		for _, cocoEntry := range cocoEntries {
			subcommand := model.NewAutocompleteData(normalizeName(cocoEntry), "[extra-text]", fmt.Sprintf("Draw the ![%s](%s =x24) coco material", normalizeName(cocoEntry), nameToUrl(cocoEntry)))
			subcommand.CanAutocompleteInTheMiddle = true
			subcommand.Replace = fmt.Sprintf("![%s](%s =x24)", normalizeName(cocoEntry), nameToUrl(cocoEntry))
			subcatcommand.AddCommand(subcommand)
		}
	}
	return coco
}
