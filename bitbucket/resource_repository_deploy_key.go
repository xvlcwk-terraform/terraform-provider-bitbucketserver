package bitbucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bitbucketTypes "github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/types"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type SSHKey struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	CreatedDate jsonTime `json:"createdDate,omitempty"`
	UpdatedDate jsonTime `json:"updatedDate,omitempty"`
}

func resourceRepositoryDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryDeployKeyCreate,
		Read:   resourceRepositoryDeployKeyRead,
		Delete: resourceRepositoryDeployKeyDelete,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"REPO_READ", "REPO_WRITE", "REPO_ADMIN"}, false),
			},
			"expiry_days": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

type KeyRequestKey struct {
	Label      string `json:"label"`
	Text       string `json:"text"`
	ExpiryDays int    `json:"expiryDays,omitempty"`
}

type KeyRequest struct {
	Key        KeyRequestKey `json:"key"`
	Permission string        `json:"permission"`
}

type KeyResponse struct {
	Key struct {
		Id          int    `json:"id"`
		ExpiryDays  int    `json:"expiryDays"`
		CreatedDate int    `json:"createdDate"`
		Label       string `json:"label"`
		Text        string `json:"text"`
	} `json:"key"`
	Permission string `json:"permission"`
	Repository struct {
		Slug    string `json:"slug"`
		Project struct {
			Key string `json:"key"`
		} `json:"project"`
	} `json:"repository"`
}

func resourceRepositoryDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	keyRequest := &KeyRequest{
		Key: KeyRequestKey{
			Label: d.Get("label").(string),
			Text:  d.Get("key").(string),
		},
		Permission: d.Get("permission").(string),
	}

	expiryDays := d.Get("expiry_days")
	if expiryDays != nil && expiryDays.(int) != 0 {
		keyRequest.Key.ExpiryDays = expiryDays.(int)
	}

	request, err := json.Marshal(keyRequest)

	client := m.(*bitbucketTypes.BitbucketServerProvider).BitbucketClient
	resp, err := client.Post(fmt.Sprintf("/rest/keys/latest/projects/%s/repos/%s/ssh",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
	), bytes.NewBuffer(request))

	if err != nil {
		return err
	}
	return storeResponse(d, resp)
}

func resourceRepositoryDeployKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*bitbucketTypes.BitbucketServerProvider).BitbucketClient

	resp, err := client.Get(fmt.Sprintf("/rest/keys/latest/projects/%s/repos/%s/ssh/%s",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
		url.QueryEscape(d.Id()),
	))
	if err != nil {
		return err
	}
	return storeResponse(d, resp)
}

func resourceRepositoryDeployKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*bitbucketTypes.BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/keys/latest/projects/%s/repos/%s/ssh/%s",
		url.QueryEscape(d.Get("project").(string)),
		url.QueryEscape(d.Get("repository").(string)),
		url.QueryEscape(d.Id()),
	))
	return err
}

func storeResponse(d *schema.ResourceData, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var keyResponse KeyResponse
	err = json.Unmarshal(body, &keyResponse)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(keyResponse.Key.Id))
	labelError := store(d, keyResponse, keyResponse.Key.Label, "label")
	permissionError := store(d, keyResponse, keyResponse.Permission, "permission")
	keyError := store(d, keyResponse, keyResponse.Key.Text, "key")
	projectError := store(d, keyResponse, keyResponse.Repository.Project.Key, "project")
	repositoryError := store(d, keyResponse, keyResponse.Repository.Slug, "repository")
	days := keyResponse.Key.ExpiryDays
	d.Set("expiry_days", days)

	return errors.Join(labelError, permissionError, keyError, projectError, repositoryError)
}

func store(d *schema.ResourceData, keyResponse KeyResponse, value string, name string) error {
	if value == "" {
		respAsJson, _ := json.Marshal(keyResponse)
		return errors.New(fmt.Sprintf("%s is nil in %s", name, string(respAsJson)))
	} else {
		return d.Set(name, value)
	}
}
