package gluster

import (
	"bytes"
	"fmt"
	"github.com/nilpntr/gluster-exporter/internal/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func execGlusterCommand(arg ...string) (*bytes.Buffer, error) {
	stdoutBuffer := &bytes.Buffer{}
	argXML := append(arg, "--xml")
	glusterExec := exec.Command(viper.GetString("gluster_binary"), argXML...)
	glusterExec.Stdout = stdoutBuffer
	err := glusterExec.Run()

	if err != nil {
		zap.L().Sugar().Errorf("tried to execute %v and got error: %v", arg, err)
		return stdoutBuffer, err
	}
	return stdoutBuffer, nil
}

func GetMountCheck() (*bytes.Buffer, error) {
	stdoutBuffer := &bytes.Buffer{}
	mountCmd := exec.Command("mount", "-t", "fuse.glusterfs")

	mountCmd.Stdout = stdoutBuffer

	return stdoutBuffer, mountCmd.Run()
}

type Mount struct {
	MountPoint string
	Volume     string
}

// ParseMountOutput pares output of system execution 'mount'
func ParseMountOutput(mountBuffer string) ([]Mount, error) {
	mounts := make([]Mount, 0, 2)
	mountRows := strings.Split(mountBuffer, "\n")
	for _, row := range mountRows {
		trimmedRow := strings.TrimSpace(row)
		if len(row) > 3 {
			mountColumns := strings.Split(trimmedRow, " ")
			mounts = append(mounts, Mount{MountPoint: mountColumns[2], Volume: mountColumns[0]})
		}
	}
	return mounts, nil
}

func ExecTouchOnVolumes(mountpoint string) (bool, error) {
	testFileName := fmt.Sprintf("%v/%v_%v", mountpoint, "gluster_mount.test", time.Now())
	_, createErr := os.Create(testFileName)
	if createErr != nil {
		return false, createErr
	}
	removeErr := os.Remove(testFileName)
	if removeErr != nil {
		return false, removeErr
	}
	return true, nil
}

// GetVolumeInfo executes "gluster volume info" at the local machine and
// returns VolumeInfoXML struct and error
func GetVolumeInfo() (VolumeInfoXML, error) {
	args := []string{"volume", "info"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return VolumeInfoXML{}, cmdErr
	}
	volumeInfo, err := utils.DecodeXml[VolumeInfoXML](bytesBuffer)
	if err != nil {
		zap.L().Sugar().Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeInfo, err
	}

	return volumeInfo, nil
}

// GetVolumeList executes "gluster volume info" at the local machine and
// returns VolumeList struct and error
func GetVolumeList() (VolList, error) {
	args := []string{"volume", "list"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return VolList{}, cmdErr
	}
	volumeList, err := utils.DecodeXml[VolumeListXML](bytesBuffer)
	if err != nil {
		zap.L().Sugar().Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeList.VolList, err
	}

	return volumeList.VolList, nil
}

// GetPeerStatus executes "gluster peer status" at the local machine and
// returns PeerStatus struct and error
func GetPeerStatus() (PeerStatus, error) {
	args := []string{"peer", "status"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return PeerStatus{}, cmdErr
	}
	peerStatus, err := utils.DecodeXml[PeerStatusXML](bytesBuffer)
	if err != nil {
		zap.L().Sugar().Errorf("Something went wrong while unmarshalling xml: %v", err)
		return peerStatus.PeerStatus, err
	}

	return peerStatus.PeerStatus, nil
}

// GetVolumeProfileGvInfoCumulative executes "gluster volume {volume] profile info cumulative" at the local machine and
// returns VolumeInfoXML struct and error
func GetVolumeProfileGvInfoCumulative(volumeName string) (VolProfile, error) {
	args := []string{"volume", "profile", volumeName, "info", "cumulative"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return VolProfile{}, cmdErr
	}
	volumeProfile, err := utils.DecodeXml[VolumeProfileXML](bytesBuffer)
	if err != nil {
		zap.L().Sugar().Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeProfile.VolProfile, err
	}
	return volumeProfile.VolProfile, nil
}

// GetVolumeStatusAllDetail executes "gluster volume status all detail" at the local machine
// returns VolumeStatusXML struct and error
func GetVolumeStatusAllDetail() (VolumeStatusXML, error) {
	args := []string{"volume", "status", "all", "detail"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return VolumeStatusXML{}, cmdErr
	}
	volumeStatus, err := utils.DecodeXml[VolumeStatusXML](bytesBuffer)
	if err != nil {
		zap.L().Sugar().Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeStatus, err
	}
	return volumeStatus, nil
}

// GetVolumeHealInfo executes volume heal info on host system and processes input
// returns (int) number of unsynced files
func GetVolumeHealInfo(volumeName string) (int, error) {
	args := []string{"volume", "heal", volumeName, "info"}
	entriesOutOfSync := 0
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return -1, cmdErr
	}
	healInfo, err := utils.DecodeXml[VolumeHealInfoXML](bytesBuffer)
	if err != nil {
		zap.L().Sugar().Error(err)
		return -1, err
	}

	for _, brick := range healInfo.HealInfo.Bricks.Brick {
		var count int
		var err error
		count, err = strconv.Atoi(brick.NumberOfEntries)
		if err != nil {
			zap.L().Sugar().Error(err)
			return -1, err
		}
		entriesOutOfSync += count
	}
	return entriesOutOfSync, nil
}

// GetVolumeQuotaList executes volume quota list on host system and processes input
// returns QuotaList structs and errors
func GetVolumeQuotaList(volumeName string) (VolumeQuotaXML, error) {
	args := []string{"volume", "quota", volumeName, "list"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return VolumeQuotaXML{}, cmdErr
	}
	volumeQuota, err := utils.DecodeXml[VolumeQuotaXML](bytesBuffer)
	if err != nil {
		zap.L().Sugar().Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeQuota, err
	}
	return volumeQuota, nil
}
