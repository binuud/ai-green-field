/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"

export enum TrainingActionTrainingActionEnum {
  TrainingActionIgnore = "TrainingActionIgnore",
  New = "New",
  Start = "Start",
  Pause = "Pause",
  Stop = "Stop",
  GetTrainingData = "GetTrainingData",
  GetTestingData = "GetTestingData",
  GetPrediction = "GetPrediction",
}

export enum ModelStateTrainingStateEnum {
  TrainingStateIgnore = "TrainingStateIgnore",
  New = "New",
  Running = "Running",
  Error = "Error",
  Pause = "Pause",
  Completed = "Completed",
  Stopped = "Stopped",
}

export enum ModelConfigModelTypeEnum {
  ModelTypeIgnore = "ModelTypeIgnore",
  Classification = "Classification",
  Regression = "Regression",
}

export enum ModelConfigActivationFunctionEnum {
  ActivationFunctionIgnore = "ActivationFunctionIgnore",
  ActivationFunctionNone = "ActivationFunctionNone",
  Relu = "Relu",
  Sigmoid = "Sigmoid",
  Tanh = "Tanh",
  Linear = "Linear",
}

export enum ModelConfigRegularizationEnum {
  RegularizationIgnore = "RegularizationIgnore",
  L1 = "L1",
  L2 = "L2",
  Dropout = "Dropout",
  RegularizationNone = "RegularizationNone",
}

export type TrainingAction = {
  action?: TrainingActionTrainingActionEnum
}

export type TrainingData = {
  x?: number[]
  y?: number[]
}

export type Prediction = {
  y?: number[]
}

export type ModelState = {
  status?: ModelStateTrainingStateEnum
  currentEpoch?: number
  trainingLoss?: number
  testLoss?: number
  createdAt?: GoogleProtobufTimestamp.Timestamp
  updatedAt?: GoogleProtobufTimestamp.Timestamp
}

export type ModelConfig = {
  name?: string
  description?: string
  learningRate?: number
  epochs?: number
  epochBatch?: number
  activationFunction?: ModelConfigActivationFunctionEnum
  regularization?: ModelConfigRegularizationEnum
  regularizationRate?: number
  trainingType?: ModelConfigModelTypeEnum
  seed?: string
  numInputs?: number
  numOutputs?: number
  numLayers?: number
}

export type Model = {
  uuid?: string
  state?: ModelState
  config?: ModelConfig
  linearModel?: LinearRegressionModel
}

export type LinearRegressionModel = {
  weight?: number
  bias?: number
}