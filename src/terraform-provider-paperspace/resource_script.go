package main

import (
  "encoding/json"
  "fmt"
  "github.com/hashicorp/terraform/helper/schema"
  "log"
)

func resourceScriptCreate(d *schema.ResourceData, m interface{}) error {
  client := m.(PaperspaceClient).RestyClient

  log.Printf("[INFO] paperspace resourceScriptCreate Client ready")

  region := m.(PaperspaceClient).Region
	if r, ok := d.GetOk("region"); ok {
		region = r.(string)
	}
  if region == "" {
    return fmt.Errorf("Error creating paperspace script: missing region")
  }

  body := make(MapIf)
  body.AppendAs(d, "name", "scriptName")
  body.Append(d, "scriptText")
  body.AppendAsIfSet(d, "description", "scriptDescription")
  body.AppendIfSet(d, "isEnabled")
  body.AppendIfSet(d, "runOnce")

  data, _ := json.MarshalIndent(body, "", "  ")
  log.Println(string(data))

  resp, err := client.R().
  SetBody(body).
  Post("/scripts/createScript")

  if err != nil {
    return fmt.Errorf("Error creating paperspace script: %s", err)
  }

  statusCode := resp.StatusCode()
  log.Printf("[INFO] paperspace resourceScriptCreate StatusCode: %v", statusCode)
  LogResponse("paperspace resourceScriptCreate", resp, err)
  if statusCode != 200 {
    return fmt.Errorf("Error creating paperspace script: Response: %s", resp.Body())
  }

  var f interface{}
  err = json.Unmarshal(resp.Body(), &f)

  if err != nil {
    return fmt.Errorf("Error unmarshalling paperspace script create response: %s", err)
  }

  mp := f.(map[string]interface{})
  id, _ := mp["id"].(string)

  if id == "" {
    return fmt.Errorf("Error in paperspace script create data: id not found")
  }

  log.Printf("[INFO] paperspace resourceScriptCreate returned id: %v", id)

  SetResData(d, mp, "name")
  SetResData(d, mp, "description")
  SetResData(d, mp, "ownerType")
  SetResData(d, mp, "ownerId")
  SetResData(d, mp, "dtCreated")
  SetResData(d, mp, "isEnabled")
  SetResData(d, mp, "runOnce")

  d.SetId(id);

  return nil
}

func resourceScriptRead(d *schema.ResourceData, m interface{}) error {
  client := m.(PaperspaceClient).RestyClient

  log.Printf("[INFO] paperspace resourceScriptRead Client ready")

  resp, err := client.R().
  Get("/scripts/getScript?scriptId=" + d.Id())

  if err != nil {
    return fmt.Errorf("Error reading paperspace script: %s", err)
  }

  statusCode := resp.StatusCode()
  log.Printf("[INFO] paperspace resourceScriptRead StatusCode: %v", statusCode)
  LogResponse("paperspace resourceScriptCreate", resp, err)
  if statusCode == 404 {
    log.Printf("[INFO] paperspace resourceScriptRead scriptId not found; removing resource %s", d.Id())
    d.SetId("")
    return nil
  }
  if statusCode != 200 {
    return fmt.Errorf("Error reading paperspace script: Response: %s", resp.Body())
  }

  var f interface{}
  err = json.Unmarshal(resp.Body(), &f)

  if err != nil {
    return fmt.Errorf("Error unmarshalling paperspace script read response: %s", err)
  }

  mp := f.(map[string]interface{})
  id, _ := mp["id"].(string)

  if id == "" {
    log.Printf("[WARNING] paperspace resourceScriptRead script id not found; removing resource %s", d.Id())
    d.SetId("")
    return nil
  }

  log.Printf("[INFO] paperspace resourceScriptRead returned id: %v", id)

  SetResData(d, mp, "name")
  SetResData(d, mp, "description")
  SetResData(d, mp, "ownerType")
  SetResData(d, mp, "ownerId")
  SetResData(d, mp, "dtCreated")
  SetResData(d, mp, "isEnabled")
  SetResData(d, mp, "runOnce")

  client = m.(PaperspaceClient).RestyClient

  resp, err = client.R().
  Get("/scripts/getScriptText?scriptId=" + d.Id())

  if err != nil {
    return fmt.Errorf("Error reading paperspace script text: %s", err)
  }

  statusCode = resp.StatusCode()
  log.Printf("[INFO] paperspace resourceScriptRead text StatusCode: %v", statusCode)
  LogResponse("paperspace resourceScriptCreate", resp, err)
  if statusCode == 404 {
    log.Printf("[INFO] paperspace resourceScriptRead text scriptId not found")
    return nil
  }
  if statusCode != 200 {
    return fmt.Errorf("Error reading paperspace script text: Response: %s", resp.Body())
  }

  d.Set("scriptText", resp.Body())

  return nil
}

func resourceScriptUpdate(d *schema.ResourceData, m interface{}) error {

  log.Printf("[INFO] paperspace resourceScriptUpdate Client ready")

  return nil
}

func resourceScriptDelete(d *schema.ResourceData, m interface{}) error {
  client := m.(PaperspaceClient).RestyClient

  log.Printf("[INFO] paperspace resourceScriptDelete Client ready")

  resp, err := client.R().
  Post("/scripts/" + d.Id() + "/destroy")

  if err != nil {
    return fmt.Errorf("Error deleting paperspace script: %s", err)
  }

  statusCode := resp.StatusCode()
  log.Printf("[INFO] paperspace resourceScriptDelete StatusCode: %v", statusCode)
  LogResponse("paperspace resourceScriptDelete", resp, err)
  if statusCode != 204 && statusCode != 404 {
    return fmt.Errorf("Error deleting paperspace script: Response: %s", resp.Body())
  }
  if statusCode == 204 {
    log.Printf("[INFO] paperspace resourceScriptDelete script deleted successfully, StatusCode: %v", statusCode)
  }
  if statusCode == 404 {
    log.Printf("[INFO] paperspace resourceScriptDelete script already deleted, StatusCode: %v", statusCode)
  }

  return nil
}

func resourceScript() *schema.Resource {
  return &schema.Resource{
    Create: resourceScriptCreate,
    Read:   resourceScriptRead,
    Update: resourceScriptUpdate,
    Delete: resourceScriptDelete,
    Importer: &schema.ResourceImporter{
        State: schema.ImportStatePassthrough,
      },

    Schema: map[string]*schema.Schema{
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
      "description": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "scriptText": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
      "ownerType": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "ownerId": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "dtCreated": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "isEnabled": &schema.Schema{
        Type:     schema.TypeBool,
        Optional: true,
      },
      "runOnce": &schema.Schema{
        Type:     schema.TypeBool,
        Optional: true,
      },
    },
  }
}
