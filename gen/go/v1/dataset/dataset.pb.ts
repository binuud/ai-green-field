/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"

export enum DatasetTrainingDataTypeEnum {
  TrainingDataTypeIgnore = "TrainingDataTypeIgnore",
  LinearEquation = "LinearEquation",
  Audio = "Audio",
  Image = "Image",
  Video = "Video",
  Text = "Text",
}

export type Dataset = {
  uuid?: string
  description?: string
  type?: DatasetTrainingDataTypeEnum
  createdAt?: GoogleProtobufTimestamp.Timestamp
  updatedAt?: GoogleProtobufTimestamp.Timestamp
}