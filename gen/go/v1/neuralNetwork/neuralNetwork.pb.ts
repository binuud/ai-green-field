/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"

export enum TrainingStateTrainingTypeEnum {
  TrainingTypeIgnore = "TrainingTypeIgnore",
  Classification = "Classification",
  Regression = "Regression",
}

export enum TrainingStateTrainingStateEnum {
  TrainingStateIgnore = "TrainingStateIgnore",
  New = "New",
  Running = "Running",
  Error = "Error",
  Pause = "Pause",
  Completed = "Completed",
  Stopped = "Stopped",
}

export enum TrainingStateActivationFunctionEnum {
  ActivationFunctionIgnore = "ActivationFunctionIgnore",
  ActivationFunctionNone = "ActivationFunctionNone",
  Relu = "Relu",
  Sigmoid = "Sigmoid",
  Tanh = "Tanh",
  Linear = "Linear",
}

export enum TrainingStateRegularizationEnum {
  RegularizationIgnore = "RegularizationIgnore",
  L1 = "L1",
  L2 = "L2",
  Dropout = "Dropout",
  RegularizationNone = "RegularizationNone",
}

export type TrainingState = {
  uuid?: string
  learningRate?: number
  epochs?: number
  epochBatch?: number
  activationFunction?: TrainingStateActivationFunctionEnum
  regularization?: TrainingStateRegularizationEnum
  regularizationRate?: number
  trainingType?: TrainingStateTrainingTypeEnum
  trainingLoss?: number
  testLoss?: number
  numInputs?: number
  numOutputs?: number
  numLayers?: number
  state?: TrainingStateTrainingStateEnum
  currentEpoch?: number
  createdAt?: GoogleProtobufTimestamp.Timestamp
  updatedAt?: GoogleProtobufTimestamp.Timestamp
}