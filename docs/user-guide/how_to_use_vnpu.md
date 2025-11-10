# HAMi ascend-device-plugin User Guide

## Introduction

Volcano supports **two vNPU modes** for sharing Ascend devices:

---

### 1. Mindcluster mode

**Description**:

MindCluster, formerly known as [Ascend/ascend-device-plugin](https://gitee.com/ascend/ascend-device-plugin) is an official device plugin, which supports npu cluster for all Ascend series and vnpu feature for Ascend 310 series. 

**Use case**:

NPU cluster for Ascend 910 series  
NPU and vNPU cluster for Ascend 310 series 

---

### 2. HAMi mode

**Description**:

This mode is developed by a third-party community 'HAMi', which is the developer of [volcano-vgpu](./how_to_use_volcano_vgpu.md) feature, It supports vNPU feature for both Ascend 310 and Ascend 910. It also support managing heterogeneous Ascend cluster(Cluster with multiple Ascend types,i.e 910A,910B2,910B3,310p)

**Use case**:

NPU and vNPU cluster for Ascend 910 series  
NPU and vNPU cluster for Ascend 310 series  
Heterogeneous Ascend cluster

---

## Installation

To enable vNPU scheduling, the following components must be set up based on the selected mode:


**Prerequisites**:

Kubernetes >= 1.16  
Volcano >= 1.14  
[ascend-docker-runtime](https://gitcode.com/Ascend/mind-cluster/tree/master/component/ascend-docker-runtime) (for HAMi Mode) 

### Install Volcano:

Follow instructions in Volcano Installer Guide

  * Follow instructions in [Volcano Installer Guide](https://github.com/volcano-sh/volcano?tab=readme-ov-file#quick-start-guide)

### Install ascend-device-plugin

In this step, you need to select different ascend-device-plugin based on the vNPU mode you selected.

---

#### MindCluster Mode

```
Wait for @JackyTYang to fill
```

---

#### HAMi mode

##### Deploy `hami-scheduler-device` config map

```
kubectl apply -f https://raw.githubusercontent.com/Project-HAMi/ascend-device-plugin/refs/heads/main/hami-scheduler-device.yaml
```

##### Deploy ascend-device-plugin

```
kubectl apply -f https://raw.githubusercontent.com/Project-HAMi/ascend-device-plugin/refs/heads/main/ascend-device-plugin.yaml
```

refer https://github.com/Project-HAMi/ascend-device-plugin

##### Scheduler Config Update
```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: volcano-scheduler-configmap
  namespace: volcano-system
data:
  volcano-scheduler.conf: |
    actions: "enqueue, allocate, backfill"
    tiers:
    - plugins:
      - name: predicates
      - name: deviceshare
        arguments:
          deviceshare.AscendHAMiVNPUEnable: true   # enable ascend vnpu
          deviceshare.SchedulePolicy: binpack  # scheduling policy. binpack / spread
          deviceshare.KnownGeometriesCMNamespace: kube-system
          deviceshare.KnownGeometriesCMName: hami-scheduler-device
```

  **Note:** You may noticed that, 'volcano-vgpu' has its own GeometriesCMName and GeometriesCMNamespace, which means if you want to use both vNPU and vGPU in a same volcano cluster, you need to merge the configMap from both sides and set it here.

## Usage

Usage is different depending on the mode you selected

---

### MindCluster mode

```
Wait for @JackyTYang to fill
```

---

### HAMi mode

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ascend-pod
spec:
  schedulerName: volcano
  containers:
    - name: ubuntu-container
      image: swr.cn-south-1.myhuaweicloud.com/ascendhub/ascend-pytorch:24.0.RC1-A2-1.11.0-ubuntu20.04
      command: ["sleep"]
      args: ["100000"]
      resources:
        limits:
          huawei.com/Ascend310P: "1"
          huawei.com/Ascend310P-memory: "4096"

```

The supported Ascend chips and their resouruceNames are shown in the following table:

| ChipName | ResourceName | ResourceMemoryName |
|-------|-------|-------|
| 910A | huawei.com/Ascend910A | huawei.com/Ascend910A-memory |
| 910B2 | huawei.com/Ascend910B2 | huawei.com/Ascend910B2-memory |
| 910B3 | huawei.com/Ascend910B3 | huawei.com/Ascend910B3-memory |
| 910B4 | huawei.com/Ascend910B4 | huawei.com/Ascend910B4-memory |
| 910B4-1 | huawei.com/Ascend910B4-1 | huawei.com/Ascend910B4-1-memory |
| 310P3 | huawei.com/Ascend310P | huawei.com/Ascend310P-memory |