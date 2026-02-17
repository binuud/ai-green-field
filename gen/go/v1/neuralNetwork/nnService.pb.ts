/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"
import * as NeuralNetworkNeuralNetwork from "./neuralNetwork.pb"
export type PingNeuralNetworkRequest = {
}

export type PingNeuralNetworkResponse = {
  createdAt?: GoogleProtobufTimestamp.Timestamp
  description?: string
  data?: string
}

export type CreateNeuralNetworkRequest = {
  state?: NeuralNetworkNeuralNetwork.TrainingState
}

export type CreateNeuralNetworkResponse = {
  state?: NeuralNetworkNeuralNetwork.TrainingState
}

export type TrainNeuralNetworkRequest = {
}

export type TrainNeuralNetworkResponse = {
}

export type TestNeuralNetworkRequest = {
}

export type TestNeuralNetworkResponse = {
}

export type LoadNeuralNetworkRequest = {
}

export type LoadNeuralNetworkResponse = {
}

export type SaveNeuralNetworkRequest = {
}

export type SaveNeuralNetworkResponse = {
}

export type ListNeuralNetworkRequest = {
}

export type ListNeuralNetworkResponse = {
}

export class NeuralNetwork {
  static Ping(req: PingNeuralNetworkRequest, initReq?: fm.InitReq): Promise<PingNeuralNetworkResponse> {
    return fm.fetchReq<PingNeuralNetworkRequest, PingNeuralNetworkResponse>(`/v1/neuralNetwork/ping?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Create(req: CreateNeuralNetworkRequest, initReq?: fm.InitReq): Promise<CreateNeuralNetworkResponse> {
    return fm.fetchReq<CreateNeuralNetworkRequest, CreateNeuralNetworkResponse>(`/v1/neuralNetwork/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Train(req: TrainNeuralNetworkRequest, initReq?: fm.InitReq): Promise<TrainNeuralNetworkResponse> {
    return fm.fetchReq<TrainNeuralNetworkRequest, TrainNeuralNetworkResponse>(`/v1/neuralNetwork/train?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Test(req: TestNeuralNetworkRequest, initReq?: fm.InitReq): Promise<TestNeuralNetworkResponse> {
    return fm.fetchReq<TestNeuralNetworkRequest, TestNeuralNetworkResponse>(`/v1/neuralNetwork/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Load(req: LoadNeuralNetworkRequest, initReq?: fm.InitReq): Promise<LoadNeuralNetworkResponse> {
    return fm.fetchReq<LoadNeuralNetworkRequest, LoadNeuralNetworkResponse>(`/v1/neuralNetwork/load`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Save(req: SaveNeuralNetworkRequest, initReq?: fm.InitReq): Promise<SaveNeuralNetworkResponse> {
    return fm.fetchReq<SaveNeuralNetworkRequest, SaveNeuralNetworkResponse>(`/v1/neuralNetwork/save`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static List(req: ListNeuralNetworkRequest, initReq?: fm.InitReq): Promise<ListNeuralNetworkResponse> {
    return fm.fetchReq<ListNeuralNetworkRequest, ListNeuralNetworkResponse>(`/v1/neuralNetwork/list?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}