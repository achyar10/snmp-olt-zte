package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/achyar10/snmp-olt-zte/config"
	"github.com/achyar10/snmp-olt-zte/internal/model"
	"github.com/achyar10/snmp-olt-zte/internal/repository"
	"github.com/achyar10/snmp-olt-zte/internal/utils"
	"github.com/gosnmp/gosnmp"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/singleflight"
)

type OnuUseCaseInterface interface {
	GetByBoardIDAndPonID(ctx context.Context, boardID, ponID int) ([]model.ONUInfoPerBoard, error)
	GetByBoardIDPonIDAndOnuID(boardID, ponID, onuID int) (model.ONUCustomerInfo, error)
	GetEmptyOnuID(ctx context.Context, boardID, ponID int) ([]model.OnuID, error)
	GetOnuIDAndSerialNumber(boardID, ponID int) ([]model.OnuSerialNumber, error)
	UpdateEmptyOnuID(ctx context.Context, boardID, ponID int) error
	GetByBoardIDAndPonIDWithPagination(boardID, ponID, page, pageSize int) (
		[]model.ONUInfoPerBoard, int,
	)
}

type onuUsecase struct {
	snmpRepository  repository.SnmpRepositoryInterface
	redisRepository repository.OnuRedisRepositoryInterface
	cfg             *config.Config
	sg              singleflight.Group
}

func NewOnuUsecase(
	snmpRepository repository.SnmpRepositoryInterface, redisRepository repository.OnuRedisRepositoryInterface,
	cfg *config.Config,
) OnuUseCaseInterface {
	return &onuUsecase{
		snmpRepository:  snmpRepository,
		redisRepository: redisRepository,
		cfg:             cfg,
		sg:              singleflight.Group{},
	}
}

// getOltInfo is a function to get OLT information
func (u *onuUsecase) getOltConfig(boardID, ponID int) (*model.OltConfig, error) {
	cfg, err := u.getBoardConfig(boardID, ponID)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return cfg, nil
}

// getBoardConfig is a function to get board configuration
func (u *onuUsecase) getBoardConfig(boardID, ponID int) (*model.OltConfig, error) {
	switch boardID {
	case 1:
		return u.getBoard1Config(ponID), nil
	case 2:
		return u.getBoard2Config(ponID), nil
	default:
		return nil, errors.New("invalid Board ID")
	}
}

// getBoard1Config is a function to get board 1 configuration
func (u *onuUsecase) getBoard1Config(ponID int) *model.OltConfig {
	// Define the configuration for Board 1
	switch ponID {
	case 1:
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon1.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon1.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon1.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon1.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon1.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon1.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon1.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon1.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon1.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon1.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon1.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon1.OnuGponOpticalDistanceOID,
		}
	case 2:
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon2.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon2.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon2.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon2.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon2.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon2.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon2.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon2.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon2.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon2.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon2.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon2.OnuGponOpticalDistanceOID,
		}
	case 3: // PON 3
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon3.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon3.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon3.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon3.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon3.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon3.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon3.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon3.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon3.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon3.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon3.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon3.OnuGponOpticalDistanceOID,
		}
	case 4: // PON 4
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon4.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon4.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon4.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon4.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon4.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon4.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon4.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon4.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon4.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon4.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon4.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon4.OnuGponOpticalDistanceOID,
		}
	case 5: // PON 5
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon5.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon5.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon5.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon5.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon5.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon5.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon5.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon5.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon5.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon5.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon5.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon5.OnuGponOpticalDistanceOID,
		}
	case 6: // PON 6
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon6.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon6.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon6.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon6.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon6.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon6.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon6.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon6.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon6.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon6.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon6.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon6.OnuGponOpticalDistanceOID,
		}
	case 7: // PON 7
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon7.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon7.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon7.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon7.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon7.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon7.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon7.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon7.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon7.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon7.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon7.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon7.OnuGponOpticalDistanceOID,
		}
	case 8: // PON 8
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon8.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon8.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon8.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon8.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon8.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon8.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon8.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon8.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon8.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon8.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon8.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon8.OnuGponOpticalDistanceOID,
		}
	case 9: // PON 9
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon9.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon9.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon9.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon9.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon9.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon9.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon9.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon9.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon9.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon9.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon9.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon9.OnuGponOpticalDistanceOID,
		}
	case 10: // PON 10
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon10.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon10.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon10.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon10.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon10.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon10.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon10.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon10.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon10.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon10.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon10.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon10.OnuGponOpticalDistanceOID,
		}
	case 11: // PON 11
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon11.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon11.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon11.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon11.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon11.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon11.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon11.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon11.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon11.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon11.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon11.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon11.OnuGponOpticalDistanceOID,
		}
	case 12: // PON 12
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon12.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon12.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon12.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon12.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon12.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon12.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon12.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon12.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon12.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon12.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon12.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon12.OnuGponOpticalDistanceOID,
		}
	case 13: // PON 13
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon13.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon13.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon13.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon13.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon13.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon13.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon13.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon13.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon13.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon13.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon13.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon13.OnuGponOpticalDistanceOID,
		}
	case 14: // PON 14
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon14.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon14.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon14.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon14.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon14.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon14.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon14.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon14.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon14.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon14.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon14.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon14.OnuGponOpticalDistanceOID,
		}
	case 15: // PON 15
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon15.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon15.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon15.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon15.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon15.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon15.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon15.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon15.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon15.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon15.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon15.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon15.OnuGponOpticalDistanceOID,
		}

	case 16: // PON 16
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board1Pon16.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board1Pon16.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board1Pon16.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board1Pon16.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board1Pon16.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board1Pon16.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board1Pon16.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board1Pon16.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board1Pon16.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board1Pon16.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board1Pon16.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board1Pon16.OnuGponOpticalDistanceOID,
		}

	default:
		log.Error().Msg("Invalid PON ID") // Log error message
		return nil
	}
}

// getBoard2Config is a function to get board 2 configuration
func (u *onuUsecase) getBoard2Config(ponID int) *model.OltConfig {
	// Define the configuration for Board 2
	switch ponID {
	case 1: // PON 1
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon1.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon1.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon1.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon1.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon1.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon1.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon1.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon1.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon1.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon1.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon1.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon1.OnuGponOpticalDistanceOID,
		}
	case 2: // PON 2
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon2.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon2.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon2.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon2.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon2.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon2.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon2.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon2.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon2.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon2.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon2.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon2.OnuGponOpticalDistanceOID,
		}
	case 3: // PON 3
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon3.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon3.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon3.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon3.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon3.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon3.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon3.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon3.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon3.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon3.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon3.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon3.OnuGponOpticalDistanceOID,
		}
	case 4: // PON 4
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon4.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon4.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon4.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon4.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon4.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon4.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon4.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon4.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon4.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon4.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon4.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon4.OnuGponOpticalDistanceOID,
		}
	case 5: // PON 5
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon5.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon5.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon5.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon5.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon5.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon5.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon5.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon5.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon5.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon5.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon5.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon5.OnuGponOpticalDistanceOID,
		}
	case 6: // PON 6
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon6.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon6.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon6.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon6.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon6.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon6.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon6.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon6.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon6.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon6.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon6.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon6.OnuGponOpticalDistanceOID,
		}
	case 7: // PON 7
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon7.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon7.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon7.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon7.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon7.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon7.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon7.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon7.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon7.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon7.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon7.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon7.OnuGponOpticalDistanceOID,
		}
	case 8: // PON 8
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon8.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon8.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon8.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon8.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon8.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon8.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon8.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon8.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon8.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon8.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon8.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon8.OnuGponOpticalDistanceOID,
		}
	case 9: // PON 9
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon9.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon9.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon9.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon9.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon9.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon9.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon9.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon9.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon9.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon9.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon9.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon9.OnuGponOpticalDistanceOID,
		}
	case 10: // PON 10
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon10.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon10.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon10.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon10.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon10.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon10.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon10.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon10.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon10.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon10.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon10.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon10.OnuGponOpticalDistanceOID,
		}
	case 11: // PON 11
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon11.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon11.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon11.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon11.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon11.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon11.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon11.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon11.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon11.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon11.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon11.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon11.OnuGponOpticalDistanceOID,
		}

	case 12: // PON 12
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon12.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon12.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon12.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon12.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon12.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon12.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon12.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon12.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon12.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon12.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon12.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon12.OnuGponOpticalDistanceOID,
		}

	case 13: // PON 13
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon13.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon13.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon13.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon13.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon13.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon13.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon13.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon13.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon13.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon13.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon13.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon13.OnuGponOpticalDistanceOID,
		}
	case 14: // PON 14
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon14.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon14.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon14.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon14.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon14.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon14.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon14.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon14.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon14.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon14.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon14.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon14.OnuGponOpticalDistanceOID,
		}
	case 15: // PON 15
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon15.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon15.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon15.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon15.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon15.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon15.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon15.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon15.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon15.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon15.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon15.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon15.OnuGponOpticalDistanceOID,
		}
	case 16: // PON 16
		return &model.OltConfig{
			BaseOID:                   u.cfg.OltCfg.BaseOID1,
			OnuIDNameOID:              u.cfg.Board2Pon16.OnuIDNameOID,
			OnuTypeOID:                u.cfg.Board2Pon16.OnuTypeOID,
			OnuSerialNumberOID:        u.cfg.Board2Pon16.OnuSerialNumberOID,
			OnuRxPowerOID:             u.cfg.Board2Pon16.OnuRxPowerOID,
			OnuTxPowerOID:             u.cfg.Board2Pon16.OnuTxPowerOID,
			OnuStatusOID:              u.cfg.Board2Pon16.OnuStatusOID,
			OnuIPAddressOID:           u.cfg.Board2Pon16.OnuIPAddressOID,
			OnuDescriptionOID:         u.cfg.Board2Pon16.OnuDescriptionOID,
			OnuLastOnlineOID:          u.cfg.Board2Pon16.OnuLastOnlineOID,
			OnuLastOfflineOID:         u.cfg.Board2Pon16.OnuLastOfflineOID,
			OnuLastOfflineReasonOID:   u.cfg.Board2Pon16.OnuLastOfflineReasonOID,
			OnuGponOpticalDistanceOID: u.cfg.Board2Pon16.OnuGponOpticalDistanceOID,
		}
	default:
		log.Error().Msg("Invalid PON ID") // Log error message
		return nil
	}
}

func (u *onuUsecase) GetByBoardIDAndPonID(ctx context.Context, boardID, ponID int) ([]model.ONUInfoPerBoard, error) {
	log.Info().Msg("Get All ONU Information from Board ID: " + strconv.Itoa(boardID) + " and PON ID: " + strconv.Itoa(ponID))

	key := fmt.Sprintf("onuinfo-b%d-p%d", boardID, ponID)

	// Using simple flight to prevent duplicate SNMP requests
	result, err, _ := u.sg.Do(key, func() (interface{}, error) {
		// Get OLT config
		oltConfig, err := u.getOltConfig(boardID, ponID) // Get OLT config based on Board ID and PON ID
		if err != nil {
			log.Error().Msg("Failed to get OLT Config: " + err.Error())
			return nil, err
		}

		// Redis key
		redisKey := fmt.Sprintf("board_%d_pon_%d", boardID, ponID)

		// Check if data is already cached in Redis
		cachedOnuData, err := u.redisRepository.GetONUInfoList(ctx, redisKey) // Get ONU Information from Redis
		if err == nil && cachedOnuData != nil {
			log.Info().Msg("Get ONU Information from Redis with Key: " + redisKey)
			return cachedOnuData, nil
		}

		// SNMP Walk to get Information from OLT Board and PON
		log.Info().Msg("Get All ONU Information from SNMP Walk Board ID: " + strconv.Itoa(boardID) + " and PON ID: " + strconv.Itoa(ponID))
		// Create a map to store SNMP Walk results
		snmpDataMap := make(map[string]gosnmp.SnmpPDU)
		// Perform SNMP Walk to get ONU ID and Name using snmpRepository Walk method with timeout context parameter
		err = u.snmpRepository.Walk(oltConfig.BaseOID+oltConfig.OnuIDNameOID, func(pdu gosnmp.SnmpPDU) error {
			snmpDataMap[utils.ExtractONUID(pdu.Name)] = pdu
			return nil
		})

		if err != nil {
			return nil, err
		}

		var onuInformationList []model.ONUInfoPerBoard // Create a slice of ONUInfoPerBoard

		// Loop through SNMP data map to get ONU information based on ONU ID and ONU Name stored in map before and store
		for _, pdu := range snmpDataMap {
			onuInfo := model.ONUInfoPerBoard{
				Board: boardID,
				PON:   ponID,
				ID:    utils.ExtractIDOnuID(pdu.Name),
				Name:  utils.ExtractName(pdu.Value),
			}

			// Get Data ONU Type from SNMP Walk using getONUType method
			if onuType, err := u.getONUType(oltConfig.OnuTypeOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.OnuType = onuType
			}
			// Get Data ONU Serial Number from SNMP Walk using getSerialNumber method
			if sn, err := u.getSerialNumber(oltConfig.OnuSerialNumberOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.SerialNumber = sn
			}
			// Get Data ONU RX Power from SNMP Walk using getRxPower method
			if rx, err := u.getRxPower(oltConfig.OnuRxPowerOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.RXPower = rx
			}
			// Get Data ONU TX Power from SNMP Walk using getTxPower method
			if status, err := u.getStatus(oltConfig.OnuStatusOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.Status = status
			}

			// Get Data ONU IP Address from SNMP Walk using getIPAddress method
			onuInformationList = append(onuInformationList, onuInfo)
		}

		// Sort the ONU information list by ID
		sort.Slice(onuInformationList, func(i, j int) bool {
			return onuInformationList[i].ID < onuInformationList[j].ID
		})

		// Save the ONU information list to Redis with a 5-minute expiration time
		err = u.redisRepository.SaveONUInfoList(ctx, redisKey, 300, onuInformationList)
		if err != nil {
			log.Error().Msg("Failed to save ONU Information to Redis: " + err.Error())
		} else {
			log.Info().Msg("Saved ONU Information to Redis with Key: " + redisKey)
		}

		// Return the ONU information list
		return onuInformationList, nil
	})

	if err != nil {
		log.Error().Msg("Failed to get ONU Information: " + err.Error()) // Log error message to logger
		return nil, err                                                  // Return error if error is not nil
	}

	return result.([]model.ONUInfoPerBoard), nil // Return the result from the cache or SNMP Walk
}

func (u *onuUsecase) GetByBoardIDPonIDAndOnuID(boardID, ponID, onuID int) (
	model.ONUCustomerInfo, error,
) {
	// Set key for simple flight
	key := fmt.Sprintf("onu:%d:%d:%d", boardID, ponID, onuID)

	// Using simple flight to prevent duplicate SNMP requests
	result, err, _ := u.sg.Do(key, func() (interface{}, error) {
		oltConfig, err := u.getOltConfig(boardID, ponID) // Get OLT config based on Board ID and PON ID
		if err != nil {
			log.Error().Msg("Failed to get OLT Config: " + err.Error())
			return model.ONUCustomerInfo{}, err
		}

		var onuInformationList model.ONUCustomerInfo   // Create a variable to store ONU information
		snmpDataMap := make(map[string]gosnmp.SnmpPDU) // Create a map to store SNMP Walk results

		log.Info().Msg("Get Detail ONU Information with SNMP Walk from Board ID: " +
			strconv.Itoa(boardID) + " PON ID: " + strconv.Itoa(ponID) +
			" ONU ID: " + strconv.Itoa(onuID))

		// Get ONU ID and Name using snmpRepository Walk method with timeout context parameter
		err = u.snmpRepository.Walk(oltConfig.BaseOID+oltConfig.OnuIDNameOID+"."+strconv.Itoa(onuID),
			func(pdu gosnmp.SnmpPDU) error {
				snmpDataMap[utils.ExtractONUID(pdu.Name)] = pdu
				return nil
			})
		if err != nil {
			log.Error().Msg("Failed to walk OID: " + err.Error())
			return model.ONUCustomerInfo{}, errors.New("failed to walk OID")
		}

		// Loop through SNMP data map to get ONU information based on ONU ID and ONU Name stored in map before and store
		for _, pdu := range snmpDataMap {
			onuInfo := model.ONUCustomerInfo{
				Board: boardID,
				PON:   ponID,
				ID:    utils.ExtractIDOnuID(pdu.Name),
				Name:  utils.ExtractName(pdu.Value),
			}

			// Get Data ONU Type from SNMP Walk using getONUType method
			if onuType, err := u.getONUType(oltConfig.OnuTypeOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.OnuType = onuType
			}

			// Get Data ONU Serial Number from SNMP Walk using getSerialNumber method
			if serial, err := u.getSerialNumber(oltConfig.OnuSerialNumberOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.SerialNumber = serial
			}

			// Get Data ONU RX Power from SNMP Walk using getRxPower method
			if rx, err := u.getRxPower(oltConfig.OnuRxPowerOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.RXPower = rx
			}

			// Get Data ONU TX Power from SNMP Walk using getTxPower method
			if tx, err := u.getTxPower(oltConfig.OnuTxPowerOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.TXPower = tx
			}

			// Get Data ONU Status from SNMP Walk using getStatus method
			if status, err := u.getStatus(oltConfig.OnuStatusOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.Status = status
			}

			// Get Data ONU IP Address from SNMP Walk using getIPAddress method
			if ip, err := u.getIPAddress(oltConfig.OnuIPAddressOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.IPAddress = ip
			}

			// Get Data ONU Description from SNMP Walk using getDescription method
			if desc, err := u.getDescription(oltConfig.OnuDescriptionOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.Description = desc
			}

			// Get Data ONU Last Online from SNMP Walk using getLastOnline method
			if lastOnline, err := u.getLastOnline(oltConfig.OnuLastOnlineOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.LastOnline = lastOnline
			}

			// Get Data ONU Last Offline from SNMP Walk using getLastOffline method
			if lastOffline, err := u.getLastOffline(oltConfig.OnuLastOfflineOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.LastOffline = lastOffline
			}

			// Get Data ONU Last Offline Reason from SNMP Walk using getLastOfflineReason method
			if uptime, err := u.getUptimeDuration(onuInfo.LastOnline); err == nil {
				onuInfo.Uptime = uptime
			}

			// Get Data ONU Last Down Time Duration from SNMP Walk using getLastDownDuration method
			if downtime, err := u.getLastDownDuration(onuInfo.LastOffline, onuInfo.LastOnline); err == nil {
				onuInfo.LastDownTimeDuration = downtime
			}

			// Get Data ONU Last Offline Reason from SNMP Walk using getLastOfflineReason method
			if reason, err := u.getLastOfflineReason(oltConfig.OnuLastOfflineReasonOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.LastOfflineReason = reason
			}

			// Get Data ONU GPON Optical Distance from SNMP Walk using getOnuGponOpticalDistance method
			if dist, err := u.getOnuGponOpticalDistance(oltConfig.OnuGponOpticalDistanceOID, strconv.Itoa(onuInfo.ID)); err == nil {
				onuInfo.GponOpticalDistance = dist
			}

			onuInformationList = onuInfo // Append ONU information to the onuInformationList
		}

		return onuInformationList, nil // Return the ONU information list
	})

	if err != nil {
		return model.ONUCustomerInfo{}, err
	}

	return result.(model.ONUCustomerInfo), nil // Return the result from the cache or SNMP Walk
}

func (u *onuUsecase) GetEmptyOnuID(ctx context.Context, boardID, ponID int) ([]model.OnuID, error) {
	// Set key for simple flight
	key := fmt.Sprintf("empty_onu_id:%d:%d", boardID, ponID)

	// Using simple flight to prevent duplicate requests for the same data
	result, err, _ := u.sg.Do(key, func() (interface{}, error) {
		// Get OLT config based on Board ID and PON ID
		oltConfig, err := u.getOltConfig(boardID, ponID)
		if err != nil {
			log.Error().Msg("Failed to get OLT Config for Get Empty ONU ID: " + err.Error())
			return nil, err
		}

		// Redis Key
		redisKey := "board_" + strconv.Itoa(boardID) + "_pon_" + strconv.Itoa(ponID) + "_empty_onu_id"

		// Try to get data from Redis using GetOnuIDCtx method with context and Redis key as parameter
		cachedOnuData, err := u.redisRepository.GetOnuIDCtx(ctx, redisKey)
		if err == nil && cachedOnuData != nil {
			log.Info().Msg("Get Empty ONU ID from Redis with Key: " + redisKey)
			// If data exists in Redis, return data from Redis
			return cachedOnuData, nil
		}

		// Perform SNMP Walk to get ONU ID and ONU Name
		snmpOID := oltConfig.BaseOID + oltConfig.OnuIDNameOID
		emptyOnuIDList := make([]model.OnuID, 0)

		log.Info().Msg("Get Empty ONU ID with SNMP Walk from Board ID: " + strconv.Itoa(boardID) + " and PON ID: " + strconv.Itoa(ponID))

		// Perform SNMP Walk to get ONU ID and Name
		err = u.snmpRepository.Walk(snmpOID, func(pdu gosnmp.SnmpPDU) error {
			idOnuID := utils.ExtractIDOnuID(pdu.Name)
			emptyOnuIDList = append(emptyOnuIDList, model.OnuID{
				Board: boardID,
				PON:   ponID,
				ID:    idOnuID,
			})
			return nil
		})
		if err != nil {
			log.Error().Msg("Failed to perform SNMP Walk get empty ONU ID: " + err.Error())
			return nil, err
		}

		// Create a map to store numbers to be deleted
		numbersToRemove := make(map[int]bool)

		for _, onuInfo := range emptyOnuIDList {
			numbersToRemove[onuInfo.ID] = true
		}

		// Remove the numbers that should not be added to the emptyOnuIDList
		emptyOnuIDList = emptyOnuIDList[:0]

		// Loop through 128 numbers to get the numbers to be deleted
		for i := 1; i <= 128; i++ {
			if _, ok := numbersToRemove[i]; !ok {
				emptyOnuIDList = append(emptyOnuIDList, model.OnuID{
					Board: boardID,
					PON:   ponID,
					ID:    i,
				})
			}
		}

		// Sort by ID ascending
		sort.Slice(emptyOnuIDList, func(i, j int) bool {
			return emptyOnuIDList[i].ID < emptyOnuIDList[j].ID
		})

		// Set data to Redis
		err = u.redisRepository.SetOnuIDCtx(ctx, redisKey, 300, emptyOnuIDList)
		if err != nil {
			log.Error().Msg("Failed to set data to Redis: " + err.Error())
			return nil, err
		}

		log.Info().Msg("Save Empty ONU ID to Redis with Key: " + redisKey)

		return emptyOnuIDList, nil
	})

	if err != nil {
		log.Error().Msg("Failed to get Empty ONU ID: " + err.Error()) // Log error message to logger
		return nil, err                                               // Return error if error is not nil
	}

	return result.([]model.OnuID), nil
}

func (u *onuUsecase) GetOnuIDAndSerialNumber(boardID, ponID int) ([]model.OnuSerialNumber, error) {
	// Set key for simple flight
	key := fmt.Sprintf("onu_id_and_serial_number:%d:%d", boardID, ponID)

	// Using simple flight to prevent duplicate requests for the same data
	result, err, _ := u.sg.Do(key, func() (interface{}, error) {
		// Get OLT config based on Board ID and PON ID
		oltConfig, err := u.getOltConfig(boardID, ponID)
		if err != nil {
			log.Error().Msg("Failed to get OLT Config: " + err.Error())
			return nil, err
		}

		// Perform SNMP Walk to get ONU ID
		snmpOID := oltConfig.BaseOID + oltConfig.OnuIDNameOID
		onuIDList := make([]model.OnuID, 0)

		log.Info().Msg("Get ONU ID with SNMP Walk from Board ID: " + strconv.Itoa(boardID) + " and PON ID: " + strconv.Itoa(ponID))

		// Perform SNMP BulkWalk to get ONU ID and Name
		err = u.snmpRepository.Walk(snmpOID, func(pdu gosnmp.SnmpPDU) error {
			idOnuID := utils.ExtractIDOnuID(pdu.Name)
			onuIDList = append(onuIDList, model.OnuID{
				Board: boardID,
				PON:   ponID,
				ID:    idOnuID,
			})
			return nil
		})
		if err != nil {
			log.Error().Msg("Failed to perform SNMP Walk get ONU ID: " + err.Error())
			return nil, err
		}

		// Create a slice of ONU Serial Number
		var onuSerialNumberList []model.OnuSerialNumber

		// Loop through onuIDList to get ONU Serial Number
		for _, onuInfo := range onuIDList {
			// Get Data ONU Serial Number from SNMP Walk using getSerialNumber method
			onuSerialNumber, err := u.getSerialNumber(oltConfig.OnuSerialNumberOID, strconv.Itoa(onuInfo.ID))
			if err == nil {
				onuSerialNumberList = append(onuSerialNumberList, model.OnuSerialNumber{
					Board:        boardID,
					PON:          ponID,
					ID:           onuInfo.ID,
					SerialNumber: onuSerialNumber,
				})
			}
		}

		// Sort ONU Serial Number list based on ONU ID ascending
		sort.Slice(onuSerialNumberList, func(i, j int) bool {
			return onuSerialNumberList[i].ID < onuSerialNumberList[j].ID
		})

		return onuSerialNumberList, nil
	})

	if err != nil {
		log.Error().Msg("Failed to get ONU ID and Serial Number: " + err.Error()) // Log error message to logger
		return nil, err                                                           // Return error if error is not nil
	}

	return result.([]model.OnuSerialNumber), nil
}

func (u *onuUsecase) UpdateEmptyOnuID(ctx context.Context, boardID, ponID int) error {
	// Set key for simple flight
	key := fmt.Sprintf("update_empty_onu_id:%d:%d", boardID, ponID)

	// Using simple flight to prevent duplicate requests for the same data
	_, err, _ := u.sg.Do(key, func() (interface{}, error) {
		// Get OLT config based on Board ID and PON ID
		oltConfig, err := u.getOltConfig(boardID, ponID)
		if err != nil {
			log.Error().Msg("Failed to get OLT Config: " + err.Error())
			return nil, err
		}

		// Perform SNMP Walk to get ONU ID and ONU Name
		snmpOID := oltConfig.BaseOID + oltConfig.OnuIDNameOID
		emptyOnuIDList := make([]model.OnuID, 0)

		log.Info().Msg("Get Empty ONU ID with SNMP Walk from Board ID: " + strconv.Itoa(boardID) + " and PON ID: " + strconv.Itoa(ponID))

		// Perform SNMP BulkWalk to get ONU ID and Name
		err = u.snmpRepository.Walk(snmpOID, func(pdu gosnmp.SnmpPDU) error {
			idOnuID := utils.ExtractIDOnuID(pdu.Name)
			emptyOnuIDList = append(emptyOnuIDList, model.OnuID{
				Board: boardID,
				PON:   ponID,
				ID:    idOnuID,
			})
			return nil
		})
		if err != nil {
			return nil, errors.New("failed to perform SNMP Walk")
		}

		// Create a map to store numbers to be deleted
		numbersToRemove := make(map[int]bool)
		for _, onuInfo := range emptyOnuIDList {
			numbersToRemove[onuInfo.ID] = true
		}

		// Filter out ONU IDs that are not empty
		emptyOnuIDList = emptyOnuIDList[:0]
		for i := 1; i <= 128; i++ {
			if _, ok := numbersToRemove[i]; !ok {
				emptyOnuIDList = append(emptyOnuIDList, model.OnuID{
					Board: boardID,
					PON:   ponID,
					ID:    i,
				})
			}
		}

		// Sort ONU IDs by ID ascending
		sort.Slice(emptyOnuIDList, func(i, j int) bool {
			return emptyOnuIDList[i].ID < emptyOnuIDList[j].ID
		})

		// Set data to Redis using SetOnuIDCtx method
		redisKey := "board_" + strconv.Itoa(boardID) + "_pon_" + strconv.Itoa(ponID) + "_empty_onu_id"
		err = u.redisRepository.SetOnuIDCtx(ctx, redisKey, 300, emptyOnuIDList)
		if err != nil {
			log.Error().Msg("Failed to set data to Redis: " + err.Error())
			return nil, errors.New("failed to set data to Redis")
		}

		log.Info().Msg("Save Update Empty ONU ID to Redis with Key: " + redisKey)
		return nil, nil
	})

	return err
}

func (u *onuUsecase) GetByBoardIDAndPonIDWithPagination(
	boardID, ponID, pageIndex, pageSize int,
) ([]model.ONUInfoPerBoard, int) {

	// Create a unique key for this request based on the parameters
	key := fmt.Sprintf("get_onu_info:%d:%d:%d:%d", boardID, ponID, pageIndex, pageSize)

	// Using simple flight to prevent duplicate requests for the same data
	result, err, _ := u.sg.Do(key, func() (interface{}, error) {
		// Get OLT config based on Board ID and PON ID
		oltConfig, err := u.getOltConfig(boardID, ponID)
		if err != nil {
			return nil, err
		}

		// SNMP OID variable
		snmpOID := oltConfig.BaseOID + oltConfig.OnuIDNameOID

		var onlyOnuIDList []model.OnuOnlyID
		var count int

		// If data does not exist in Redis, then get data from SNMP
		if len(onlyOnuIDList) == 0 {
			err := u.snmpRepository.Walk(snmpOID, func(pdu gosnmp.SnmpPDU) error {
				onlyOnuIDList = append(onlyOnuIDList, model.OnuOnlyID{
					ID: utils.ExtractIDOnuID(pdu.Name),
				})
				return nil
			})

			if err != nil {
				return nil, err
			}
		} else {
			// Optionally, handle Redis case here
			log.Error().Msg("Failed to get data from Redis")
		}

		// Calculate total count
		count = len(onlyOnuIDList)

		// Calculate the index of the first item to be retrieved
		startIndex := (pageIndex - 1) * pageSize

		// Calculate the index of the last item to be retrieved
		endIndex := startIndex + pageSize

		// If the index of the last item to be retrieved is greater than the number of items, set it to the number of items
		if endIndex > len(onlyOnuIDList) {
			endIndex = len(onlyOnuIDList)
		}

		// Slice the data for pagination
		onlyOnuIDList = onlyOnuIDList[startIndex:endIndex]

		var onuInformationList []model.ONUInfoPerBoard

		// Loop through onlyOnuIDList to get ONU information based on ONU ID
		for _, onuID := range onlyOnuIDList {
			onuInfo := model.ONUInfoPerBoard{
				Board: boardID,  // Set Board ID to ONUInfo struct Board field
				PON:   ponID,    // Set PON ID to ONUInfo struct PON field
				ID:    onuID.ID, // Set ONU ID to ONUInfo struct ID field
			}

			// Get Name based on ONU ID and ONU Name OID and store it to ONU onuInfo struct
			onuName, err := u.getName(oltConfig.OnuIDNameOID, strconv.Itoa(onuInfo.ID))
			if err == nil {
				onuInfo.Name = onuName // Set ONU Name to ONU onuInfo struct Name field
			}

			// Get ONU Type based on ONU ID and ONU Type OID and store it to ONU onuInfo struct
			onuType, err := u.getONUType(oltConfig.OnuTypeOID, strconv.Itoa(onuInfo.ID))
			if err == nil {
				onuInfo.OnuType = onuType // Set ONU Type to ONU onuInfo struct OnuType field
			}

			// Get ONU Serial Number based on ONU ID and ONU Serial Number OID and store it to ONU onuInfo struct
			onuSerialNumber, err := u.getSerialNumber(oltConfig.OnuSerialNumberOID, strconv.Itoa(onuInfo.ID))
			if err == nil {
				onuInfo.SerialNumber = onuSerialNumber // Set ONU Serial Number to ONU onuInfo struct SerialNumber field
			}

			// Get ONU RX Power based on ONU ID and ONU RX Power OID and store it to ONU onuInfo struct
			onuRXPower, err := u.getRxPower(oltConfig.OnuRxPowerOID, strconv.Itoa(onuInfo.ID))
			if err == nil {
				onuInfo.RXPower = onuRXPower // Set ONU RX Power to ONU onuInfo struct RXPower field
			}

			// Get ONU Status based on ONU ID and ONU Status OID and store it to ONU onuInfo struct
			onuStatus, err := u.getStatus(oltConfig.OnuStatusOID, strconv.Itoa(onuInfo.ID))
			if err == nil {
				onuInfo.Status = onuStatus // Set ONU Status to ONU onuInfo struct Status field
			}

			// Append ONU information to the onuInformationList
			onuInformationList = append(onuInformationList, onuInfo)
		}

		// Sort ONU information list based on ONU ID ascending
		sort.Slice(onuInformationList, func(i, j int) bool {
			return onuInformationList[i].ID < onuInformationList[j].ID
		})

		// Return both the list and the count inside a struct
		return model.PaginationResult{
			OnuInformationList: onuInformationList,
			Count:              count,
		}, nil
	})

	// Handle error if any occurred during simple flight processing
	if err != nil {
		return nil, 0
	}

	// Extract the result from the simpleflight result and return it
	paginationResult := result.(model.PaginationResult)
	return paginationResult.OnuInformationList, paginationResult.Count

}

func (u *onuUsecase) getName(OnuIDNameOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID1 + OnuIDNameOID + "." + onuID
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}
	return utils.ExtractName(result.Variables[0].Value), nil
}

func (u *onuUsecase) getONUType(OnuTypeOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID2 + OnuTypeOID + "." + onuID
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}
	return utils.ExtractName(result.Variables[0].Value), nil
}

func (u *onuUsecase) getSerialNumber(OnuSerialNumberOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID1 + OnuSerialNumberOID + "." + onuID
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}
	return utils.ExtractSerialNumber(result.Variables[0].Value), nil
}

func (u *onuUsecase) getTxPower(OnuTxPowerOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID2 + OnuTxPowerOID + "." + onuID + ".1"
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}
	power, _ := utils.ConvertAndMultiply(result.Variables[0].Value)
	return power, nil
}

func (u *onuUsecase) getRxPower(OnuRxPowerOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID1 + OnuRxPowerOID + "." + onuID + ".1"
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}
	power, _ := utils.ConvertAndMultiply(result.Variables[0].Value)
	return power, nil
}

func (u *onuUsecase) getStatus(OnuStatusOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID1 + OnuStatusOID + "." + onuID
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}
	return utils.ExtractAndGetStatus(result.Variables[0].Value), nil
}

func (u *onuUsecase) getIPAddress(OnuIPAddressOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID2 + OnuIPAddressOID + "." + onuID + ".1"
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}
	return utils.ExtractName(result.Variables[0].Value), nil
}

func (u *onuUsecase) getDescription(OnuDescriptionOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID1 + OnuDescriptionOID + "." + onuID
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}
	return utils.ExtractName(result.Variables[0].Value), nil
}

func (u *onuUsecase) getLastOnline(OnuLastOnlineOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID1 + OnuLastOnlineOID + "." + onuID
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}

	value := result.Variables[0].Value.([]byte)
	return utils.ConvertByteArrayToDateTime(value)
}

func (u *onuUsecase) getLastOffline(OnuLastOfflineOID, onuID string) (string, error) {
	baseOID := u.cfg.OltCfg.BaseOID1
	oid := baseOID + OnuLastOfflineOID + "." + onuID
	oids := []string{oid}

	result, err, _ := u.sg.Do(oid, func() (interface{}, error) {
		return u.snmpRepository.Get(oids)
	})
	if err != nil {
		log.Error().Msg("Failed to perform SNMP Get for last offline: " + err.Error())
		return "", errors.New("failed to perform SNMP Get")
	}

	resultData := result.(*gosnmp.SnmpPacket)
	if len(resultData.Variables) > 0 {
		value := resultData.Variables[0].Value.([]byte)
		return utils.ConvertByteArrayToDateTime(value)
	}

	log.Error().Msg("Failed to get ONU Last Offline: No variables in the response")
	return "", errors.New("no variables in the response")
}

func (u *onuUsecase) getLastOfflineReason(OnuLastOfflineReasonOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID1 + OnuLastOfflineReasonOID + "." + onuID
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}

	return utils.ExtractLastOfflineReason(result.Variables[0].Value), nil
}

func (u *onuUsecase) getOnuGponOpticalDistance(OnuGponOpticalDistanceOID, onuID string) (string, error) {
	oid := u.cfg.OltCfg.BaseOID1 + OnuGponOpticalDistanceOID + "." + onuID
	result, err := u.getFromSNMPWithSingleflight(oid)
	if err != nil {
		return "", err
	}

	return utils.ExtractGponOpticalDistance(result.Variables[0].Value), nil
}

func (u *onuUsecase) getUptimeDuration(lastOnline string) (string, error) {
	currentTime := time.Now()

	lastOnlineTime, err := time.Parse("2006-01-02 15:04:05", lastOnline)
	if err != nil {
		log.Error().Msg("Failed to parse last online time: " + err.Error())
		return "", err
	}

	duration := currentTime.Sub(lastOnlineTime) + time.Hour*7
	return utils.ConvertDurationToString(duration), nil
}

// Last Down Duration
func (u *onuUsecase) getLastDownDuration(lastOffline, lastOnline string) (string, error) {
	lastOfflineTime, err := time.Parse("2006-01-02 15:04:05", lastOffline)
	if err != nil {
		log.Error().Msg("Failed to parse last offline time: " + err.Error())
		return "", err
	}

	lastOnlineTime, err := time.Parse("2006-01-02 15:04:05", lastOnline)
	if err != nil {
		log.Error().Msg("Failed to parse last online time: " + err.Error())
		return "", err
	}

	duration := lastOnlineTime.Sub(lastOfflineTime)
	return utils.ConvertDurationToString(duration), nil
}

func (u *onuUsecase) getFromSNMPWithSingleflight(oid string) (*gosnmp.SnmpPacket, error) {
	result, err, _ := u.sg.Do(oid, func() (interface{}, error) {
		return u.snmpRepository.Get([]string{oid})
	})
	if err != nil {
		log.Error().Msg("Failed to perform SNMP Get for OID " + oid + ": " + err.Error())
		return nil, errors.New("failed to perform SNMP Get")
	}

	packet := result.(*gosnmp.SnmpPacket)
	if len(packet.Variables) == 0 {
		log.Error().Msg("No variables returned for OID " + oid)
		return nil, errors.New("no variables in the response")
	}

	return packet, nil
}
