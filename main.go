package main

import (
	"fmt"
	"strconv"
	"time"
  "go.uber.org/multierr"
)

type FieldType int

const (
  Date FieldType = iota // time.Time
  Category // string
  Number // float64
)

type Field interface {
  String() string  
  FromString(s string) error 
}

type FieldDef struct {
  Name string
  Type FieldType  
}

func NewFieldDef(tp FieldType, name string) *FieldDef {
  f:= new(FieldDef)
  f.Name=name
  f.Type=tp
  return f
}

type CategoryField struct {
  Data string
}

func (c CategoryField) String() string {
  return c.Data
}

func (c *CategoryField) FromString(s string) error {
  c.Data = s
  return nil
}

type NumberField struct {
  Data float64
}

func (c NumberField) String() string {
  return fmt.Sprint(c.Data)
}

func (c *NumberField) FromString(s string) error {
  f64, err := strconv.ParseFloat(s, 64)
  if err == nil {
	 c.Data =  f64
  } 
  return err
}

type DateField struct {
  Data time.Time
}

func (c DateField) String() string {
  return fmt.Sprint(c.Data)
}

func (c *DateField) FromString(s string) error{
  d, err := time.Parse(dateformat, s)
  if err == nil {
	 c.Data =  d
  } 
  return err
}

type Row []Field

func (r Row)String() string {
  s:=""
  for i, field := range r  {
    s+=field.String()
    if i< (len(r)-1) {
      s+=", "
    }
  }  
  return s 
}

type DataRow struct {
  Definitions []FieldDef
  Fields []Row  
  Errors []error
  throwError bool
}

func NewDataRow(throwError bool) *DataRow {
  dr := new(DataRow)
  dr.throwError=throwError
  return dr
}

func (d *DataRow) AddField(f FieldDef) {
  d.Definitions=append(d.Definitions,f)  
}


func (d *DataRow) Set(values ...string) (err error) {   
  line:=len(d.Fields)-1
  
  if len(values) > len(d.Definitions) {
    err = fmt.Errorf("length of array [%d] does not match the amount of values [%d]",len(d.Definitions),len(values))
  } else {
    // create the cells in the row
    var row Row
    for _,def := range d.Definitions {
      switch(def.Type) {
        case Category:
          row = append(row,new(CategoryField))
            case Date:
          row = append(row,new(DateField))
            case Number:
          row = append(row,new(NumberField))
      }    
    }
    // add values to the row
    for i, val := range values {
      err = multierr.Append(err, row[i].FromString(val))    
    }

    // add row into data frame
    d.Fields=append(d.Fields,row)
    line++
  }
  if d.throwError && (err !=nil) {
    return fmt.Errorf("[%d]-%v",line,err)
  } 
  if !d.throwError && (err !=nil) {
    d.Errors = append(d.Errors,fmt.Errorf("[%d]-%v",line,err))  
  }
  return nil
}

func (d DataRow) String() string {
  s:=""
  for j,row := range d.Fields {
    s+=fmt.Sprint(j)+", "+ row.String() + "\n"
  }
  
  return s
}

func (d DataRow) Errorf() string {
  s:=""
  for _,row := range d.Errors {
    s+=fmt.Sprint(row) + "\n"
  } 
  return s
}

var dateformat ="02-01-2006 15:04:05"

func SetDateFormat(ds string) {
  dateformat = ds
}

func TestDateFormat(ds string) string {
  d, err := time.Parse(dateformat, ds)
  if err != nil {
	 return err.Error()
  } 
  return d.Format(dateformat)
}

func main() {
  dr := NewDataRow(false) 
  dr.AddField(*NewFieldDef(Date,"created_at"))
  dr.AddField(*NewFieldDef(Category,"product_category"))
  dr.AddField(*NewFieldDef(Number,"cost"))

  dr.Set("15-07-1990 10:22:04","Hobs","1.35")
  dr.Set("23-12-2024X23:59:04","","123341.35")
  dr.Set("15-07-1990 10:22:04","Hobs","2.67")
  dr.Set("23-12-2024 23:59:04","Horno","123341;35")
  dr.Set("15-07-2023 10:22:04","Hobs","8.88")
 
  fmt.Println(dr)
  fmt.Println(dr.Errorf())
}
