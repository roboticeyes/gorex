// Copyright 2019 Robotic Eyes. All rights reserved.

package core

import (
	"fmt"
)

// ProjectDetailedList can be used to get a detailed list of project owned by somebody
type ProjectDetailedList struct {
	Embedded struct {
		Projects []struct {
			Public                     bool   `json:"public"`
			NumberOfProjectFiles       int    `json:"numberOfProjectFiles"`
			TotalProjectFileSize       int    `json:"totalProjectFileSize"`
			RootRexReferenceKey        string `json:"rootRexReferenceKey"`
			NumberOfReadPermittedUsers int    `json:"numberOfReadPermittedUsers"`
			LastUpdated                string `json:"lastUpdated"`
			DateCreated                string `json:"dateCreated"`
			Owner                      string `json:"owner"`
			Name                       string `json:"name"`
			Urn                        string `json:"urn"`
			Links                      struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				Project struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"project"`
				RexReferences struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"rexReferences"`
				RootRexReference struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"rootRexReference"`
				ProjectFiles struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"projectFiles"`
			} `json:"_links"`
		} `json:"projects"`
	} `json:"_embedded"`
}

// Project is the structure of a REX project
type Project struct {
	DateCreated string `json:"dateCreated"`
	CreatedBy   string `json:"createdBy"`
	LastUpdated string `json:"lastUpdated"`
	UpdatedBy   string `json:"updatedBy"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	TagLine     string `json:"tagLine"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Urn         string `json:"urn"`
	Embedded    struct {
		RootRexReference struct {
			RootReference bool   `json:"rootReference"`
			Key           string `json:"key"`
			Links         struct {
				Self struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"self"`
				Project struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"project"`
				ParentReference struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"parentReference"`
				ChildReferences struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"childReferences"`
				ProjectFiles struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"projectFiles"`
			} `json:"_links"`
		} `json:"rootRexReference"`
		ProjectFiles []struct {
			LastModified string `json:"lastModified"`
			FileSize     int    `json:"fileSize"`
			Name         string `json:"name"`
			Type         string `json:"type"`
			Links        struct {
				Self struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"self"`
				RexReference struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"rexReference"`
				Project struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"project"`
				FileDownload struct {
					Href string `json:"href"`
				} `json:"file.download"`
			} `json:"_links"`
		} `json:"projectFiles"`
		RexReferences []struct {
			RootReference bool   `json:"rootReference"`
			Key           string `json:"key"`
			Links         struct {
				Self struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"self"`
				Project struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"project"`
				ParentReference struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"parentReference"`
				ChildReferences struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"childReferences"`
				ProjectFiles struct {
					Href      string `json:"href"`
					Templated bool   `json:"templated"`
				} `json:"projectFiles"`
			} `json:"_links"`
		} `json:"rexReferences"`
	} `json:"_embedded"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Project struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"project"`
		ThumbnailUpload struct {
			Href string `json:"href"`
		} `json:"thumbnail.upload"`
		ThumbnailDownload struct {
			Href string `json:"href"`
		} `json:"thumbnail.download"`
		ProjectFavorite struct {
			Href string `json:"href"`
		} `json:"projectFavorite"`
		RootRexReference struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"rootRexReference"`
		ProjectFiles struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"projectFiles"`
		ProjectAcls struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"projectAcls"`
		RexReferences struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"rexReferences"`
	} `json:"_links"`
}

// String nicely prints a project
func (p Project) String() string {

	hasRootRef := false
	if p.Embedded.RootRexReference.RootReference {
		hasRootRef = true
	}

	s := fmt.Sprintf("|------------------------------------------------------------------------------------------|\n")
	s += fmt.Sprintf("| Name           | %-71s |\n", p.Name)
	s += fmt.Sprintf("| Urn            | %-71s |\n", p.Urn)
	s += fmt.Sprintf("| Owner          | %-71s |\n", p.Owner)
	s += fmt.Sprintf("| Type           | %-71s |\n", p.Type)
	s += fmt.Sprintf("| Has root ref   | %-71t |\n", hasRootRef)
	s += fmt.Sprintf("| Total files    | %-71d |\n", len(p.Embedded.ProjectFiles))
	s += fmt.Sprintf("| Total refs     | %-71d |\n", len(p.Embedded.RexReferences))

	sz := 0
	for _, f := range p.Embedded.ProjectFiles {
		sz += f.FileSize
	}
	s += fmt.Sprintf("| Total size (KB)| %-71d |\n", sz/1024)
	s += fmt.Sprintf("|------------------------------------------------------------------------------------------|\n")

	for i, f := range p.Embedded.ProjectFiles {
		length := min(35, len(f.Name))
		s += fmt.Sprintf("| %3d | %-35s | %8d (kb) | %s |\n", i, f.Name[0:length], f.FileSize/1024, f.LastModified)
	}
	s += fmt.Sprintf("|------------------------------------------------------------------------------------------|\n")

	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
