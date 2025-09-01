package service

import (
	"rustdesk-api/config"
	"rustdesk-api/lib/jwt"
	"rustdesk-api/lib/lock"
	"rustdesk-api/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service struct {
	//AdminService     *AdminService
	//AdminRoleService *AdminRoleService
	*UserService
	*AddressBookService
	*TagService
	*PeerService
	*GroupService
	*OauthService
	*LoginLogService
	*AuditService
	*ShareRecordService
	*ServerCmdService
	*LdapService
	*AppService
}

type Dependencies struct {
	Config *config.Config
	DB     *gorm.DB
	Logger *log.Logger
	Jwt    *jwt.Jwt
	Lock   *lock.Locker
}

var Config *config.Config
var DB *gorm.DB
var Logger *log.Logger
var Jwt *jwt.Jwt
var Lock lock.Locker

var AllService *Service

func New(c *config.Config, g *gorm.DB, l *log.Logger, j *jwt.Jwt, lo lock.Locker) *Service {
	Config = c
	DB = g
	Logger = l
	Jwt = j
	Lock = lo
	AllService = new(Service)
	return AllService
}

func Paginate(page, pageSize uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		if pageSize == 0 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

func CommonEnable() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", model.COMMON_STATUS_ENABLE)
	}
}
