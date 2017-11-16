package data

import (
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/genofire/golang-lib/database"
	"github.com/genofire/warehost/lib"
)

// Login found
type Login struct {
	ID          int64
	Username    string    `gorm:"type:varchar(255);unique;column:mail" json:"username"`
	Password    string    `gorm:"type:varchar(255);column:password" json:"-"`
	Active      bool      `gorm:"default:false;column:active" json:"active"`
	Code        string    `gorm:"type:varchar(255);column:code" json:"-"`
	Superadmin  bool      `gorm:"default:false;column:superadmin" json:"superadmin"`
	CreateAt    time.Time `sql:"default:current_timestamp" gorm:"column:createat" json:"createat"`
	LastLoginAt time.Time `gorm:"column:lastloginat" json:"lastloginat"`
	Invites     []Invite  `gorm:"foreignkey:Login" json:"invites"`
}

// GetInvitedby of current login -> Invitor
func (l *Login) GetInvitedby(dbconnection *gorm.DB) (invited Invite) {
	invited = Invite{InvitedID: l.ID}
	dbconnection.Where("invited = ?", invited.InvitedID).First(&invited)
	return
}

// Invite struct
type Invite struct {
	LoginID   int64 `sql:"type:bigint REFERENCES login(id) ON UPDATE CASCADE ON DELETE CASCADE;column:login;primary_key"`
	Login     Login `gorm:"column:login" json:"login"`
	InvitedID int64 `sql:"type:bigint REFERENCES login(id) ON UPDATE CASCADE ON DELETE CASCADE;column:invited;primary_key"`
	Invited   Login `gorm:"column:invited" json:"invited"`
	Admin     bool  `sql:"default:false" json:"admin"`
}

// Profil struct
type Profil struct {
	ID       int64
	LoginID  int64  `sql:"type:bigint NOT NULL UNIQUE REFERENCES login(id) ON UPDATE CASCADE ON DELETE CASCADE;column:login" json:"-"`
	Login    *Login `gorm:"foreignkey:Login;" json:"login"`
	Reseller bool   `sql:"default:false;column:reseller" json:"reseller"`
}

// TableName of struct
func (Profil) TableName() string { return "host_profil" }

// Domain struct
type Domain struct {
	ID       int64
	ProfilID int64   `sql:"type:bigint NOT NULL REFERENCES host_profil(id) ON UPDATE CASCADE ON DELETE CASCADE;column:profil" json:"-"`
	Profil   *Profil `gorm:"foreignkey:Profil;" json:"profil"`
	FQDN     string  `sql:"type:varchar(255) NOT NULL UNIQUE;column:fqdn"  json:"fqdn"`
	Code     string  `sql:"type:varchar(255);column:code"  json:"-"`
	Active   bool    `sql:"default:false;column:active" json:"active"`
	Mail     bool    `sql:"default:false;column:mail" json:"mail"`
	Web      bool    `sql:"default:false;column:web" json:"web"`
}

// TableName of struct
func (Domain) TableName() string { return "host_domain" }

// Mail struct
type Mail struct {
	ID       int64
	DomainID int64             `sql:"type:bigint NOT NULL REFERENCES host_domain(id) ON UPDATE CASCADE ON DELETE CASCADE;column:domain" json:"-"`
	Domain   *Domain           `gorm:"foreignkey:Domain;unique_index:idx_host_domain_mail" json:"domain"`
	Name     string            `sql:"type:varchar(255);column:name" gorm:"unique_index:idx_host_domain_mail" json:"name"`
	Forwards []*MailForward    `json:"forwards"`
	LoginID  lib.JSONNullInt64 `sql:"type:bigint NULL REFERENCES login(id) ON UPDATE CASCADE ON DELETE CASCADE;column:login" json:"login"`
}

// TableName of struct
func (Mail) TableName() string { return "host_mail" }

// MailForward is a Object on with address a copy of the mail should be send
type MailForward struct {
	ID     int64
	MailID int64  `sql:"type:bigint NOT NULL REFERENCES host_mail(id) ON UPDATE CASCADE ON DELETE CASCADE;column:mail" json:"-"`
	Mail   *Mail  `gorm:"foreignkey:Mail;unique_index:idx_host_domain_mail_forward" json:"mail"`
	To     string `sql:"type:varchar(255);column:to" gorm:"unique_index:idx_host_domain_mail_forward" json:"to"`
}

// TableName of struct
func (MailForward) TableName() string { return "host_mail_forward" }

// SyncModels to verify the database schema
func init() {
	database.AddModel(&Login{})
	database.AddModel(&Invite{})
	database.AddModel(&Profil{})
	database.AddModel(&Domain{})
	database.AddModel(&Mail{})
	database.AddModel(&MailForward{})
}
func CreateDatabase() {
	var result int64
	database.Read.Model(&Login{}).Count(&result)

	if result <= 0 {
		login := &Login{
			Username:   "root",
			Active:     true,
			Superadmin: true,
			Password:   lib.NewHash("root"),
		}

		database.Write.Create(login)
		log.Error("have to create \"login\"")
	} else {
		log.Info("Connection to \"login\" works")
	}
}
