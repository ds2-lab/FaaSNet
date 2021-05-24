# Alibaba Cloud Function Compute Function Cold Start Traces

## Introduction

Our trace dataset is a subset of the data described in, and analyzed, in our ATC '21 paper. Traces were obtained by
collecting 24-hour production-level traces from two datacenters
in [Alibaba Function Compute](https://www.alibabacloud.com/product/function-compute) service during May 2021.

## Using the data

### License

The data is made available and licensed under
a [Apache License 2.0](https://github.com/mason-leap-lab/FaaSNet/blob/main/LICENSE). By downloading it or using them,
you agree to the terms of this license.

### Attribution

If you use this data for a publication or project, please cite the accompanying paper:
> Ao Wang, Shuai Chang, Huangshi Tian, Hongqi Wang, Haoran Yang, Huiba Li, Rui Du, Yue Cheng. "[FAASNET: Scalable and Fast Provisioning of Custom Serverless ContainerRuntimes at Alibaba Cloud Function Compute](https://www.usenix.org/conference/atc21/presentation/wang-ao)", in Proceedings of the 2021 USENIX Annual Technical Conference (USENIX ATC 21). USENIX Association, July 2020.

Lastly, if you have any questions, comments, or concerns, or if you would like to share tools for working with the
traces, please contact us at [**awang24@gmu.edu**](mailto:awang24@gmu.edu)

### Downloading
You can download the dateset here: [LINK](https://drive.google.com/file/d/1YLkLhbeYwxobfMtY_5LWQZyHR_ewg6HK/view?usp=sharing
)

## Schema and Description

Field | Description
 :--- | :---
`__time__` | TimeStamp in seconds
`functionName` | Unique id for the function name
`latency` | Function cold start latency<sup>1</sup> in seconds
`runtime` | Function runtime (Python, nodejs, custom-runtime, etc)
`memoryMB` | Function's allocate memory in MB

Notes:

1. The function cold start latency only counts the system level's latency, such as container initialization etc, instead
   of end-to-end cold start latency.

