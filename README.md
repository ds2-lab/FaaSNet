# FaaSNet

FaaSNet is the first system that provides an end-to-end, integrated solution for FaaS-optimized container runtime provisioning. FaaSNet uses lightweight, decentralized, and adaptive Function Trees (FTs) to avoid major platform bottlenecks. 

Our USENIX ATC'21 paper: [FaaSNet: Scalable and Fast Provisioning of Custom Serverless Container Runtimes at Alibaba Cloud Function Compute](https://www.usenix.org/conference/atc21/presentation/wang-ao)

Download the preprint version on [arXiv](https://arxiv.org/abs/2105.11229). 

This repo contains two components:

- A FaaS cold start [trace dataset](#dataset-downloading) collected from [Alibaba Cloud Function Compute](https://www.alibabacloud.com/product/function-compute).
- The source code of the [function tree prototype](#function-tree-prototype) of FaaSNet.


## The Function Cold Start Traces from Alibaba Cloud Function Compute

### Introduction

Our trace dataset is a subset of the data described in, and analyzed, in our ATC '21 paper. The traces were obtained by collecting 24-hour production-level logs from two datacenters (Beijing and Shanghai) in [Alibaba Function Compute](https://www.alibabacloud.com/product/function-compute) service during May 2021.

### Using the Data

#### License

The data is made available and licensed under an [Apache License 2.0](https://github.com/mason-leap-lab/FaaSNet/blob/main/LICENSE). By downloading it or using them, you agree to the terms of this license.



#### Dataset Downloading

We would like you could help us take the survey before downloading the FaaS cold start traces. The link is as following [Google Form](https://forms.gle/pMucV9LsZXGymrr59), you could retrieve the download link after submmting the form successfully. 

Please free to let us know if you have any further questions.

### Schema and Description

Field | Description
 :--- | :---
`__time__` | TimeStamp in seconds
`functionName` | Unique id for the function name
`latency` | Function cold start latency<sup>1</sup> in seconds
`runtime` | Function runtime (Python, nodejs, custom-runtime, etc)
`memoryMB` | Function's allocated memory in MB 

Notes:

1. The function cold start latency only counts the system level's latency, such as container initialization, etc, instead
   of end-to-end cold start latency.


## Function Tree Prototype

Our released function tree (FT) [prototype](https://github.com/mason-leap-lab/FaaSNet/tree/main/functionTree) is the version that we evaluated in the ATC '21 paper. We are continuing to improve the performance of it. We're happy to accept contributions! Please feel free to hack on the FT and integrate it into your framework/platform :-).


## Attribution

If you use our trace dataset and/or the FT prototype for a publication or project, please cite the accompanying paper using this [bibtex](#to-cite-faasnet):
> Ao Wang, Shuai Chang, Huangshi Tian, Hongqi Wang, Haoran Yang, Huiba Li, Rui Du, Yue Cheng. "[FaaSNet: Scalable and Fast Provisioning of Custom Serverless ContainerRuntimes at Alibaba Cloud Function Compute](https://www.usenix.org/conference/atc21/presentation/wang-ao)", in Proceedings of the 2021 USENIX Annual Technical Conference (USENIX ATC 21). USENIX Association, July 2020.

Lastly, if you have any questions, comments, or concerns, or if you would like to share tools for working with the traces, please contact us at [**awang24@gmu.edu**](mailto:awang24@gmu.edu).


### To Cite FaaSNet

```
@inproceedings {273798,
author = {Ao Wang and Shuai Chang and Huangshi Tian and Hongqi Wang and Haoran Yang and Huiba Li and Rui Du and Yue Cheng},
title = {FaaSNet: Scalable and Fast Provisioning of Custom Serverless Container Runtimes at Alibaba Cloud Function Compute},
booktitle = {2021 {USENIX} Annual Technical Conference ({USENIX} {ATC} 21)},
year = {2021},
isbn = {978-1-939133-23-6},
pages = {443--457},
url = {https://www.usenix.org/conference/atc21/presentation/wang-ao},
publisher = {{USENIX} Association},
month = jul,
}
```
