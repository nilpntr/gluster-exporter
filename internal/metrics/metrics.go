package metrics

import (
	"errors"
	"github.com/nilpntr/gluster-exporter/internal/gluster"
	"github.com/nilpntr/gluster-exporter/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"slices"
	"strings"
)

const (
	namespace  = "gluster"
	allVolumes = "_all"
)

type Metrics struct {
	hostname string
	volumes  []string

	up                     *prometheus.Desc
	volumesCount           *prometheus.Desc
	volumeStatus           *prometheus.Desc
	nodeSizeFreeBytes      *prometheus.Desc
	nodeSizeTotalBytes     *prometheus.Desc
	nodeInodesTotal        *prometheus.Desc
	nodeInodesFree         *prometheus.Desc
	brickCount             *prometheus.Desc
	brickDuration          *prometheus.Desc
	brickDataRead          *prometheus.Desc
	brickDataWritten       *prometheus.Desc
	brickFopHits           *prometheus.Desc
	brickFopLatencyAvg     *prometheus.Desc
	brickFopLatencyMin     *prometheus.Desc
	brickFopLatencyMax     *prometheus.Desc
	peersConnected         *prometheus.Desc
	healInfoFilesCount     *prometheus.Desc
	volumeWriteable        *prometheus.Desc
	mountSuccessful        *prometheus.Desc
	quotaHardLimit         *prometheus.Desc
	quotaSoftLimit         *prometheus.Desc
	quotaUsed              *prometheus.Desc
	quotaAvailable         *prometheus.Desc
	quotaSoftLimitExceeded *prometheus.Desc
	quotaHardLimitExceeded *prometheus.Desc
}

func New() (*Metrics, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	if !utils.FileExists(viper.GetString("gluster_binary")) {
		return nil, errors.New("gluster binary not found")
	}

	volumes := strings.Split(viper.GetString("gluster_volumes"), ",")
	if len(volumes) < 1 {
		return nil, errors.New("no gluster volumes provided")
	}

	var (
		up = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Was the last query of Gluster successful.",
			nil, nil,
		)

		volumesCount = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "volumes_available"),
			"How many volumes were up at the last query.",
			nil, nil,
		)

		volumeStatus = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "volume_status"),
			"Status code of requested volume.",
			[]string{"volume"}, nil,
		)

		nodeSizeFreeBytes = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "node_size_free_bytes"),
			"Free bytes reported for each node on each instance. Labels are to distinguish origins",
			[]string{"hostname", "path", "volume"}, nil,
		)

		nodeSizeTotalBytes = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "node_size_bytes_total"),
			"Total bytes reported for each node on each instance. Labels are to distinguish origins",
			[]string{"hostname", "path", "volume"}, nil,
		)

		nodeInodesTotal = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "node_inodes_total"),
			"Total inodes reported for each node on each instance. Labels are to distinguish origins",
			[]string{"hostname", "path", "volume"}, nil,
		)

		nodeInodesFree = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "node_inodes_free"),
			"Free inodes reported for each node on each instance. Labels are to distinguish origins",
			[]string{"hostname", "path", "volume"}, nil,
		)

		brickCount = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "brick_available"),
			"Number of bricks available at last query.",
			[]string{"volume"}, nil,
		)

		brickDuration = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "brick_duration_seconds_total"),
			"Time running volume brick in seconds.",
			[]string{"volume", "brick"}, nil,
		)

		brickDataRead = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "brick_data_read_bytes_total"),
			"Total amount of bytes of data read by brick.",
			[]string{"volume", "brick"}, nil,
		)

		brickDataWritten = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "brick_data_written_bytes_total"),
			"Total amount of bytes of data written by brick.",
			[]string{"volume", "brick"}, nil,
		)

		brickFopHits = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "brick_fop_hits_total"),
			"Total amount of file operation hits.",
			[]string{"volume", "brick", "fop_name"}, nil,
		)

		brickFopLatencyAvg = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "brick_fop_latency_avg"),
			"Average fileoperations latency over total uptime",
			[]string{"volume", "brick", "fop_name"}, nil,
		)

		brickFopLatencyMin = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "brick_fop_latency_min"),
			"Minimum fileoperations latency over total uptime",
			[]string{"volume", "brick", "fop_name"}, nil,
		)

		brickFopLatencyMax = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "brick_fop_latency_max"),
			"Maximum fileoperations latency over total uptime",
			[]string{"volume", "brick", "fop_name"}, nil,
		)

		peersConnected = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "peers_connected"),
			"Is peer connected to gluster cluster.",
			nil, nil,
		)

		healInfoFilesCount = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "heal_info_files_count"),
			"File count of files out of sync, when calling 'gluster v heal VOLNAME info",
			[]string{"volume"}, nil)

		volumeWriteable = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "volume_writeable"),
			"Writes and deletes file in Volume and checks if it is writeable",
			[]string{"volume", "mountpoint"}, nil)

		mountSuccessful = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "mount_successful"),
			"Checks if mountpoint exists, returns a bool value 0 or 1",
			[]string{"volume", "mountpoint"}, nil)

		quotaHardLimit = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "volume_quota_hardlimit"),
			"Quota hard limit (bytes) in a volume",
			[]string{"path", "volume"}, nil)

		quotaSoftLimit = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "volume_quota_softlimit"),
			"Quota soft limit (bytes) in a volume",
			[]string{"path", "volume"}, nil)

		quotaUsed = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "volume_quota_used"),
			"Current data (bytes) used in a quota",
			[]string{"path", "volume"}, nil)

		quotaAvailable = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "volume_quota_available"),
			"Current data (bytes) available in a quota",
			[]string{"path", "volume"}, nil)

		quotaSoftLimitExceeded = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "volume_quota_softlimit_exceeded"),
			"Is the quota soft-limit exceeded",
			[]string{"path", "volume"}, nil)

		quotaHardLimitExceeded = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "volume_quota_hardlimit_exceeded"),
			"Is the quota hard-limit exceeded",
			[]string{"path", "volume"}, nil)
	)

	return &Metrics{
		hostname:               hostname,
		volumes:                volumes,
		up:                     up,
		volumesCount:           volumesCount,
		volumeStatus:           volumeStatus,
		nodeSizeFreeBytes:      nodeSizeFreeBytes,
		nodeSizeTotalBytes:     nodeSizeTotalBytes,
		nodeInodesTotal:        nodeInodesTotal,
		nodeInodesFree:         nodeInodesFree,
		brickCount:             brickCount,
		brickDuration:          brickDuration,
		brickDataRead:          brickDataRead,
		brickDataWritten:       brickDataWritten,
		brickFopHits:           brickFopHits,
		brickFopLatencyAvg:     brickFopLatencyAvg,
		brickFopLatencyMin:     brickFopLatencyMin,
		brickFopLatencyMax:     brickFopLatencyMax,
		peersConnected:         peersConnected,
		healInfoFilesCount:     healInfoFilesCount,
		volumeWriteable:        volumeWriteable,
		mountSuccessful:        mountSuccessful,
		quotaHardLimit:         quotaHardLimit,
		quotaSoftLimit:         quotaSoftLimit,
		quotaUsed:              quotaUsed,
		quotaAvailable:         quotaAvailable,
		quotaSoftLimitExceeded: quotaSoftLimitExceeded,
		quotaHardLimitExceeded: quotaHardLimitExceeded,
	}, nil
}

func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.up
	ch <- m.volumeStatus
	ch <- m.volumesCount
	ch <- m.brickCount
	ch <- m.brickDuration
	ch <- m.brickDataRead
	ch <- m.brickDataWritten
	ch <- m.peersConnected
	ch <- m.nodeSizeFreeBytes
	ch <- m.nodeSizeTotalBytes
	ch <- m.brickFopHits
	ch <- m.brickFopLatencyAvg
	ch <- m.brickFopLatencyMin
	ch <- m.brickFopLatencyMax
	ch <- m.healInfoFilesCount
	ch <- m.volumeWriteable
	ch <- m.mountSuccessful
	ch <- m.quotaHardLimit
	ch <- m.quotaSoftLimit
	ch <- m.quotaUsed
	ch <- m.quotaAvailable
	ch <- m.quotaSoftLimitExceeded
	ch <- m.quotaHardLimitExceeded
}

func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	// Collect metrics from volume info
	volumeInfo, err := gluster.GetVolumeInfo()
	// Couldn't parse xml, so something is really wrong and up=0
	if err != nil {
		zap.L().Sugar().Errorf("couldn't parse xml volume info: %v", err)
		ch <- prometheus.MustNewConstMetric(
			m.up, prometheus.GaugeValue, 0.0,
		)
	}

	// use OpErrno as indicator for up
	if volumeInfo.OpErrno != 0 {
		ch <- prometheus.MustNewConstMetric(
			m.up, prometheus.GaugeValue, 0.0,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			m.up, prometheus.GaugeValue, 1.0,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		m.volumesCount, prometheus.GaugeValue, float64(volumeInfo.VolInfo.Volumes.Count),
	)

	for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
		if m.volumes[0] == allVolumes || slices.Contains(m.volumes, volume.Name) {

			ch <- prometheus.MustNewConstMetric(
				m.brickCount, prometheus.GaugeValue, float64(volume.BrickCount), volume.Name,
			)

			ch <- prometheus.MustNewConstMetric(
				m.volumeStatus, prometheus.GaugeValue, float64(volume.Status), volume.Name,
			)
		}
	}

	// reads gluster peer status
	peerStatus, peerStatusErr := gluster.GetPeerStatus()
	if peerStatusErr != nil {
		zap.L().Sugar().Errorf("couldn't parse xml of peer status: %v", peerStatusErr)
	}
	count := 0
	for range peerStatus.Peer {
		count++
	}
	ch <- prometheus.MustNewConstMetric(
		m.peersConnected, prometheus.GaugeValue, float64(count),
	)

	// reads profile info
	if viper.GetBool("profile") {
		for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
			if m.volumes[0] == allVolumes || slices.Contains(m.volumes, volume.Name) {
				volumeProfile, execVolProfileErr := gluster.GetVolumeProfileGvInfoCumulative(volume.Name)
				if execVolProfileErr != nil {
					zap.L().Sugar().Errorf("Error while executing or marshalling gluster profile output: %v", execVolProfileErr)
				}
				for _, brick := range volumeProfile.Brick {
					if strings.HasPrefix(brick.BrickName, m.hostname) {
						ch <- prometheus.MustNewConstMetric(
							m.brickDuration, prometheus.CounterValue, float64(brick.CumulativeStats.Duration), volume.Name, brick.BrickName,
						)

						ch <- prometheus.MustNewConstMetric(
							m.brickDataRead, prometheus.CounterValue, float64(brick.CumulativeStats.TotalRead), volume.Name, brick.BrickName,
						)

						ch <- prometheus.MustNewConstMetric(
							m.brickDataWritten, prometheus.CounterValue, float64(brick.CumulativeStats.TotalWrite), volume.Name, brick.BrickName,
						)
						for _, fop := range brick.CumulativeStats.FopStats.Fop {
							ch <- prometheus.MustNewConstMetric(
								m.brickFopHits, prometheus.CounterValue, float64(fop.Hits), volume.Name, brick.BrickName, fop.Name,
							)

							ch <- prometheus.MustNewConstMetric(
								m.brickFopLatencyAvg, prometheus.GaugeValue, fop.AvgLatency, volume.Name, brick.BrickName, fop.Name,
							)

							ch <- prometheus.MustNewConstMetric(
								m.brickFopLatencyMin, prometheus.GaugeValue, fop.MinLatency, volume.Name, brick.BrickName, fop.Name,
							)

							ch <- prometheus.MustNewConstMetric(
								m.brickFopLatencyMax, prometheus.GaugeValue, fop.MaxLatency, volume.Name, brick.BrickName, fop.Name,
							)
						}
					}
				}
			}
		}
	}

	// executes gluster status all detail
	volumeStatusAll, err := gluster.GetVolumeStatusAllDetail()
	if err != nil {
		zap.L().Sugar().Errorf("couldn't parse xml of peer status: %v", err)
	}
	for _, vol := range volumeStatusAll.VolStatus.Volumes.Volume {
		for _, node := range vol.Node {
			ch <- prometheus.MustNewConstMetric(
				m.nodeSizeTotalBytes, prometheus.CounterValue, float64(node.SizeTotal), node.Hostname, node.Path, vol.VolName,
			)

			ch <- prometheus.MustNewConstMetric(
				m.nodeSizeFreeBytes, prometheus.GaugeValue, float64(node.SizeFree), node.Hostname, node.Path, vol.VolName,
			)
			ch <- prometheus.MustNewConstMetric(
				m.nodeInodesTotal, prometheus.CounterValue, float64(node.InodesTotal), node.Hostname, node.Path, vol.VolName,
			)

			ch <- prometheus.MustNewConstMetric(
				m.nodeInodesFree, prometheus.GaugeValue, float64(node.InodesFree), node.Hostname, node.Path, vol.VolName,
			)
		}
	}
	vols := m.volumes
	if vols[0] == allVolumes {
		zap.L().Sugar().Warn("no Volumes were given.")
		volumeList, volumeListErr := gluster.GetVolumeList()
		if volumeListErr != nil {
			zap.L().Sugar().Error(volumeListErr)
		}
		vols = volumeList.Volume
	}

	for _, vol := range vols {
		filesCount, volumeHealErr := gluster.GetVolumeHealInfo(vol)
		if volumeHealErr == nil {
			ch <- prometheus.MustNewConstMetric(
				m.healInfoFilesCount, prometheus.CounterValue, float64(filesCount), vol,
			)
		}
	}

	mountBuffer, execMountCheckErr := gluster.GetMountCheck()
	if execMountCheckErr != nil {
		zap.L().Sugar().Error(execMountCheckErr)
	} else {
		mounts, err := gluster.ParseMountOutput(mountBuffer.String())
		if err != nil {
			zap.L().Sugar().Error(err)
			if len(mounts) > 0 {
				for _, mount := range mounts {
					ch <- prometheus.MustNewConstMetric(
						m.mountSuccessful, prometheus.GaugeValue, float64(0), mount.Volume, mount.MountPoint,
					)
				}
			}
		} else {
			for _, mount := range mounts {
				ch <- prometheus.MustNewConstMetric(
					m.mountSuccessful, prometheus.GaugeValue, float64(1), mount.Volume, mount.MountPoint,
				)

				isWriteable, err := gluster.ExecTouchOnVolumes(mount.MountPoint)
				if err != nil {
					zap.L().Sugar().Error(err)
				}
				if isWriteable {
					ch <- prometheus.MustNewConstMetric(
						m.volumeWriteable, prometheus.GaugeValue, float64(1), mount.Volume, mount.MountPoint,
					)
				} else {
					ch <- prometheus.MustNewConstMetric(
						m.volumeWriteable, prometheus.GaugeValue, float64(0), mount.Volume, mount.MountPoint,
					)
				}
			}
		}
	}
	if viper.GetBool("quota") {
		for _, volume := range volumeInfo.VolInfo.Volumes.Volume {
			if m.volumes[0] == allVolumes || slices.Contains(m.volumes, volume.Name) {
				volumeQuotaXML, err := gluster.GetVolumeQuotaList(volume.Name)
				if err != nil {
					zap.L().Sugar().Error("Cannot create quota metrics if quotas are not enabled in your gluster server")
				} else {
					for _, limit := range volumeQuotaXML.VolQuota.QuotaLimits {
						ch <- prometheus.MustNewConstMetric(
							m.quotaHardLimit,
							prometheus.CounterValue,
							float64(limit.HardLimit),
							limit.Path,
							volume.Name,
						)

						ch <- prometheus.MustNewConstMetric(
							m.quotaSoftLimit,
							prometheus.CounterValue,
							float64(limit.SoftLimitValue),
							limit.Path,
							volume.Name,
						)
						ch <- prometheus.MustNewConstMetric(
							m.quotaUsed,
							prometheus.CounterValue,
							float64(limit.UsedSpace),
							limit.Path,
							volume.Name,
						)

						ch <- prometheus.MustNewConstMetric(
							m.quotaAvailable,
							prometheus.CounterValue,
							float64(limit.AvailSpace),
							limit.Path,
							volume.Name,
						)

						slExceeded := 0.0
						if limit.SlExceeded != "No" {
							slExceeded = 1.0
						}
						ch <- prometheus.MustNewConstMetric(
							m.quotaSoftLimitExceeded,
							prometheus.CounterValue,
							slExceeded,
							limit.Path,
							volume.Name,
						)

						hlExceeded := 0.0
						if limit.HlExceeded != "No" {
							hlExceeded = 1.0
						}
						ch <- prometheus.MustNewConstMetric(
							m.quotaHardLimitExceeded,
							prometheus.CounterValue,
							hlExceeded,
							limit.Path,
							volume.Name,
						)
					}
				}
			}
		}
	}
}
