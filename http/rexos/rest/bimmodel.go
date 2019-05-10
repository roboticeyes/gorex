// Copyright 2019 Robotic Eyes. All rights reserved.

package rest

// BimModel is the main BIM model structure
type BimModel struct {
	DateCreated     interface{} `json:"dateCreated"`
	CreatedBy       interface{} `json:"createdBy"`
	LastUpdated     string      `json:"lastUpdated"`
	UpdatedBy       string      `json:"updatedBy"`
	Name            string      `json:"name"`
	GlobalID        string      `json:"globalId"`
	Description     interface{} `json:"description"`
	RexDataBlockIds interface{} `json:"rexDataBlockIds"`
	BimProperties   []struct {
		Name      string `json:"name"`
		Type      string `json:"type"`
		ValueText string `json:"valueText"`
	} `json:"bimProperties"`
	Owner    string `json:"owner"`
	Urn      string `json:"urn"`
	Embedded struct {
		BimSites []struct {
			RexDataBlockIds []int  `json:"rexDataBlockIds"`
			GlobalID        string `json:"globalId"`
			Name            string `json:"name"`
			Links           struct {
				Self struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"self"`
				BimModel struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"bimModel"`
				BimBuildings struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"bimBuildings"`
			} `json:"_links"`
		} `json:"bimSites"`
	} `json:"_embedded"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		BimModel struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"bimModel"`
		SpatialStructure struct {
			Href string `json:"href"`
		} `json:"spatialStructure"`
		RexFile struct {
			Href string `json:"href"`
		} `json:"rexFile"`
		IfcFileUploadCityBim struct {
			Href string `json:"href"`
		} `json:"ifcFile.uploadCityBim"`
		BimSites struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"bimSites"`
	} `json:"_links"`
}
