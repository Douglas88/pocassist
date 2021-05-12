package database

type Webapp struct {
	Id       	int    		`gorm:"primary_key" json:"id"`
	Name		string		`gorm:"column:name" json:"name"`
	Provider	string		`gorm:"column:provider" json:"provider"`
	Remarks		string		`gorm:"column:remarks" json:"remarks"`
}

type Vulnerability struct {
	Id       	int         `gorm:"primary_key" json:"id"`
	NameZh 		string  `gorm:"column:name_zh" json:"name_zh"`
	Cve 		string     `gorm:"column:cve" json:"cve"`
	Cnnvd 		string   `gorm:"column:cnnvd" json:"cnnvd"`
	Severity	string    `gorm:"column:severity" json:"severity"`
	Category 	string    `gorm:"column:category" json:"category"`
	Description string    `gorm:"type:longtext" json:"description"`
	Suggestion	string  `gorm:"type:longtext" json:"suggestion"`
	Language 	string    `gorm:"column:language" json:"language"`
	Webapp		int      `gorm:"column:webapp" json:"webapp"`
	ForeignWebapp 		*Webapp `gorm:"foreignKey:Webapp"`
}

type VulnerabilitySearchField struct {
	Search       string
	CategoryField  string
	WebappField int
}

func GetVulnerabilitiesTotal(field *VulnerabilitySearchField) (total int64){
	db := GlobalDB.Model(&Vulnerability{})

	if field.CategoryField != ""{
		db = db.Where("category = ?", field.CategoryField)
	}
	if field.WebappField != -1{
		db = db.Where("webapp = ?", field.WebappField)
	}
	if field.Search != ""{
		db = db.Where(
			GlobalDB.Where("name_zh like ?", "%"+field.Search+"%").
				Or("cve like ?", "%"+field.Search+"%").
				Or("cnnvd like ?", "%"+field.Search+"%").
				Or("description like ?", "%"+field.Search+"%"))
	}
	db.Count(&total)
	return
}

func GetVulnerabilities(page int, pageSize int, field *VulnerabilitySearchField) (vuls []Vulnerability) {

	db := GlobalDB.Preload("ForeignWebapp")

	if field.CategoryField != ""{
		db = db.Where("category = ?", field.CategoryField)
	}
	if field.WebappField != -1{
		db = db.Where("webapp = ?", field.WebappField)
	}
	if field.Search != ""{
		db = db.Where(
			GlobalDB.Where("name_zh like ?", "%"+field.Search+"%").
				Or("cve like ?", "%"+field.Search+"%").
				Or("cnnvd like ?", "%"+field.Search+"%").
				Or("description like ?", "%"+field.Search+"%"))
	}
	//	分页
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&vuls)
	}
	return
}

func GetVulnerability(id int) (vul Vulnerability){
	GlobalDB.Model(&Vulnerability{}).Where("id = ?", id).First(&vul)
	return
}

func EditVulnerability(id int, vul Vulnerability) bool {
	GlobalDB.Model(&Vulnerability{}).Model(&Vulnerability{}).Where("id = ?", id).Updates(vul)
	return true
}

func AddVulnerability(vul Vulnerability) bool {
	GlobalDB.Create(&vul)
	return true
}


func DeleteVulnerability(id int) bool {
	GlobalDB.Model(&Vulnerability{}).Where("id = ?", id).Delete(Vulnerability{})
	return true
}

func ExistVulnerabilityByID(id int) bool {
	var vul Vulnerability
	GlobalDB.Model(&Vulnerability{}).Where("id = ?", id).First(&vul)
	if vul.Id >0 {
		return true
	}
	return false
}

func ExistVulnerabilityByNameZh(name_zh string) bool {
	var vul Vulnerability
	GlobalDB.Model(&Vulnerability{}).Where("name_zh = ?", name_zh).First(&vul)
	if vul.Id >0 {
		return true
	}
	return false
}

func GetWebAppsTotal() (total int64) {
	db := GlobalDB.Model(&Webapp{})
	db.Count(&total)
	return
}

func GetWebApps(page int, pageSize int) (apps []Webapp) {
	//	分页
	db := GlobalDB.Model(&Webapp{})
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&apps)
	}
	return
}

func AddWebapp(app Webapp) bool {
	GlobalDB.Create(&app)
	return true
}

func ExistWebappByName(name string) bool {
	var app Webapp
	GlobalDB.Model(&Webapp{}).Where("name = ?", name).First(&app)
	if app.Id >0 {
		return true
	}
	return false
}