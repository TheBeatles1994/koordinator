/*
Copyright 2022 The Koordinator Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/koordinator-sh/koordinator/pkg/koordlet/util/system"
	"github.com/koordinator-sh/koordinator/pkg/util"
)

// NOTE: functions in this file can be overwritten for extension

// GetPodCgroupDirWithKube gets the full pod cgroup parent dir with the podParentDir (excluding kubepods dir).
// @podKubeRelativeDir kubepods-burstable.slice/kubepods-pod7712555c_ce62_454a_9e18_9ff0217b8941.slice/
// @return kubepods.slice/kubepods-burstable.slice/kubepods-pod7712555c_ce62_454a_9e18_9ff0217b8941.slice/
func GetPodCgroupDirWithKube(podParentDir string) string {
	return filepath.Join(system.CgroupPathFormatter.ParentDir, podParentDir)
}

// GetPodKubeRelativePath gets the full pod cgroup parent with the pod info.
// @return like kubepods-burstable.slice/kubepods-pod7712555c_ce62_454a_9e18_9ff0217b8941.slice/
func GetPodKubeRelativePath(pod *corev1.Pod) string {
	qosClass := util.GetKubeQosClass(pod)
	return filepath.Join(
		system.CgroupPathFormatter.QOSDirFn(qosClass),
		system.CgroupPathFormatter.PodDirFn(qosClass, string(pod.UID)),
	)
}

func GetKubeQoSByCgroupParent(cgroupDir string) corev1.PodQOSClass {
	if strings.Contains(cgroupDir, "besteffort") {
		return corev1.PodQOSBestEffort
	} else if strings.Contains(cgroupDir, "burstable") {
		return corev1.PodQOSBurstable
	}
	return corev1.PodQOSGuaranteed
}

// @return like kubepods.slice/kubepods-burstable.slice/
func GetPodQoSRelativePath(qosClass corev1.PodQOSClass) string {
	return filepath.Join(
		system.CgroupPathFormatter.ParentDir,
		system.CgroupPathFormatter.QOSDirFn(qosClass),
	)
}

// GetContainerCgroupPathWithKube gets the full container cgroup parent dir with the podKubeRelativeDir and the
// containerStatus.
// @parentDir kubepods-burstable.slice/kubepods-pod7712555c_ce62_454a_9e18_9ff0217b8941.slice/
// @return kubepods.slice/kubepods-burstable.slice/kubepods-pod7712555c_ce62_454a_9e18_9ff0217b8941.slice/****.scope
func GetContainerCgroupPathWithKube(podParentDir string, c *corev1.ContainerStatus) (string, error) {
	return GetContainerCgroupPathWithKubeByID(podParentDir, c.ContainerID)
}

// GetContainerCgroupPathWithKubeByID gets the full container cgroup parent dir with the podKubeRelativeDir and the
// container ID.
// @parentDir kubepods-burstable.slice/kubepods-pod7712555c_ce62_454a_9e18_9ff0217b8941.slice/
// @return kubepods.slice/kubepods-burstable.slice/kubepods-pod7712555c_ce62_454a_9e18_9ff0217b8941.slice/****.scope
func GetContainerCgroupPathWithKubeByID(podParentDir string, containerID string) (string, error) {
	containerDir, err := system.CgroupPathFormatter.ContainerDirFn(containerID)
	if err != nil {
		return "", err
	}
	return filepath.Join(
		GetPodCgroupDirWithKube(podParentDir),
		containerDir,
	), nil
}
