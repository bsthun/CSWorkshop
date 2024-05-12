package preview

import (
	"github.com/gofiber/fiber/v2"

	"backend/common"
	"backend/service/preview"
	"backend/type/payload"
	"backend/type/response"
	"backend/type/table"
)

func GetStateHandler(c *fiber.Ctx) error {
	// * Query component
	var components []*table.SystemComponent
	if err := cc.DB.Preload("SideAGroup").Preload("SideBGroup").Find(&components).Error; err != nil {
		return response.Error(true, "Unable to query component", err)
	}

	sideAMap := make(map[string]*payload.ComponentInfo)
	sideBMap := make(map[string]*payload.ComponentInfo)
	for _, component := range components {
		if component.SideAGroup != nil {
			componentInfo, err := preview.CreateInfo(component, component.SideAGroup)
			if err != nil {
				return err
			}
			sideAMap[*component.Name] = componentInfo
		}
		if component.SideBGroup != nil {
			componentInfo, err := preview.CreateInfo(component, component.SideBGroup)
			if err != nil {
				return err
			}
			sideBMap[*component.Name] = componentInfo
		}
	}

	// * Response
	state := &payload.StateResponse{
		SideA: sideAMap,
		SideB: sideBMap,
	}

	// * Return
	return c.JSON(response.Info(state))
}
