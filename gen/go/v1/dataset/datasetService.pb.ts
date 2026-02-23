/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"
import * as DatasetDataset from "./dataset.pb"
export type PingDatasetRequest = {
}

export type PingDatasetResponse = {
  createdAt?: GoogleProtobufTimestamp.Timestamp
  description?: string
  data?: string
}

export type GetDatasetRequest = {
}

export type GetDatasetResponse = {
  dataset?: DatasetDataset.Dataset
}

export type ListDatasetRequest = {
}

export type ListDatasetResponse = {
  datasets?: DatasetDataset.Dataset[]
}

export class DatasetService {
  static Ping(req: PingDatasetRequest, initReq?: fm.InitReq): Promise<PingDatasetResponse> {
    return fm.fetchReq<PingDatasetRequest, PingDatasetResponse>(`/v1/dataset/ping?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Get(req: GetDatasetRequest, initReq?: fm.InitReq): Promise<GetDatasetResponse> {
    return fm.fetchReq<GetDatasetRequest, GetDatasetResponse>(`/v1/dataset/get?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static List(req: ListDatasetRequest, initReq?: fm.InitReq): Promise<ListDatasetResponse> {
    return fm.fetchReq<ListDatasetRequest, ListDatasetResponse>(`/v1/dataset/list?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}