# rokid-openvoice
## 开放平台接口定义文档

> 文档版本:V0.2

Rokid openvoice开放服务包含以下四部分功能，

- 设备认证
- 语音转文字（ASR）
- 语义理解（NLU）
- 自然语言合成（TTS）

## 文档简介

此文档用于定义开放平台上云端应用接口开发协议，协议遵循 [grpc](http://www.grpc.io/docs/) 协议。

## 认证

采用grpc的[调用证书](http://www.grpc.io/docs/guides/auth.html)方案，约定如下：

| 参数             | 类型     | 描述                     | 默认值  |
| -------------- | ------ | ---------------------- | ---- |
| key            | string | 开放接口Key,在管理平台获取        | 无，必填 |
| device_type_id | string | 设备类型ID                 | 无，必填 |
| device_id      | string | 设备ID                   | 无，必填 |
| service        | string | asr,nlp,tts            | 无，必填 |
| version        | string | 接口版本号                  | 无，必填 |
| time           | string | unix时间戳                | 无，必填 |
| sign           | string | 由以上几项+secret按约定的加密方式生成 | 无，必填 |

sign的生成加密方式：

key={key}&device_type_id={device_type_id}&device_id={device_id}&service={service}&version={version}&time={time}&secret={secret}

的utf8字符串的md5值

其中{xxx}由xxx的值替代

key及secret由开发方通过管理平台获取，并保管。

## Speech 接口

### protobuf 定义

```

syntax = "proto3";

package rokid.open;

service Speech {
  rpc auth(AuthRequest) returns (AuthResponse) { }

  rpc asr(stream AsrRequest) returns (stream AsrResponse) { }

  rpc nlp(NlpRequest) returns (NlpResponse) { }

  rpc tts(TtsRequest) returns (stream TtsResponse) { }

  rpc speechv(stream VoiceSpeechRequest) returns (stream SpeechResponse) { }

  rpc speecht(TextSpeechRequest) returns (SpeechResponse) { }
}

message AuthRequest {
  string key              = 1;
  string device_type_id   = 2;
  string device_id        = 3;
  string service          = 4;
  string version          = 5;
  string timestamp        = 6;
  string sign             = 7;
}

message AuthResponse {
  // 0: success
  // -1: failed
  int32 result            = 1;
}

message AsrRequest {
  oneof request_content {
    AsrHeader header      = 1;
    bytes voice           = 2;
  }
}

message AsrHeader {
  int32 id                = 1;
  string lang             = 2;
  string codec            = 3;
  // vt = voice trigger
  string vt               = 4;
}

message AsrResponse {
  string asr              = 1;
}

message NlpRequest {
  NlpHeader header        = 1;
  string asr              = 2;
}

message NlpHeader {
  int32 id              = 1;
  string lang           = 2;
  string cdomain        = 3;
}

message NlpResponse {
  string nlp              = 1;
}

message TtsRequest {
  TtsHeader header        = 1;
  string text             = 2;
}

message TtsHeader {
  int32 id              = 1;
  string declaimer      = 2;
  string codec          = 3;
}

message TtsResponse {
  string text             = 1;
  bytes voice             = 2;
}

message VoiceSpeechRequest {
  oneof request_content {
    SpeechHeader header = 1;

    bytes voice         = 2;
  }
}

message SpeechResponse {
  string asr              = 1;

  string nlp              = 2;

  string action           = 3;
}

message TextSpeechRequest {
  SpeechHeader header = 1;

  string asr          = 2;
}

message SpeechHeader {
  int32 id = 1;

  // zh
  // en
  string lang = 2;

  // pcm
  // opu
  // opu2
  string codec = 3;

  // vt = voice trigger
  string vt = 4;

  // stack of current domains
  string cdomain = 5;

  // json format
  string device = 6;
}
```

### AuthRequest

在同一条连接上第一个请求必须为Auth()。

| 参数             | 类型     | 描述                     | 默认值  |
| -------------- | ------ | ---------------------- | ---- |
| key            | string | 开放接口Key,在管理平台获取        | 无，必填 |
| device_type_id | string | 设备类型ID                 | 无，必填 |
| device_id      | string | 设备ID                   | 无，必填 |
| service        | string | asr,nlp,tts            | 无，必填 |
| version        | string | 接口版本号                  | 无，必填 |
| time           | string | unix时间戳                | 无，必填 |
| sign           | string | 由以上几项+secret按约定的加密方式生成 | 无，必填 |

sign的生成加密方式：

key={key}&device_type_id={device_type_id}&device_id={device_id}&service={service}&version={version}&time={time}&secret={secret}

的utf8字符串的md5值

其中{xxx}由xxx的值替代

key及secret由开发方通过管理平台获取，并保管。

### AuthResponse

| 参数     | 类型    | 描述                |
| ------ | ----- | ----------------- |
| result | int32 | 认证结果(成功:0, 失败:-1) |

### AsrRequest

一个请求一个AsrRequest{AsrHeader}后跟多个AsrRequest{voice}。PCM为
RIFF (little-endian) data, WAVE audio, Microsoft PCM, 16 bit, mono 16000 Hz语音。

| 参数     | 类型        | 描述                | 默认值  |
| ------ | --------- | ----------------- | ---- |
| header | AsrHeader | 帮助识别voice语音流的参数配置 | 无    |
| voice  | bytes     | 需要识别的voice语音流     | 无    |

### AsrHeader

| 参数    | 类型     | 描述                        | 默认值   |
| ----- | ------ | ------------------------- | ----- |
| id    | int32  | 唯一标识，用于跟踪一个完整的请求，处理及响应事件。 | 0     |
| lang  | string | 语音流的语言，目前支持zh-CN，en-US。   | zh-CN |
| codec | string | 语音流的编码，目前支持PCM，OPU，OPU2。  | PCM   |
| vt    | string | 激活词，即用于唤醒设备的名字，如"若琪"。     | 空     |

### AsrResponse

| 参数   | 类型     | 描述         |
| ---- | ------ | ---------- |
| asr  | string | asr实时识别的结果 |

### NlpRequest

| 参数     | 类型        | 描述                     | 默认值  |
| ------ | --------- | ---------------------- | ---- |
| header | NlpHeader | 帮助识别voice语音流的NlpHeader | 无    |
| asr    | string    | 需要识别的asr文本             | 无    |

### NlpHeader

| 参数      | 类型     | 描述                        | 默认值   |
| ------- | ------ | ------------------------- | ----- |
| id      | int32  | 唯一标识，用于跟踪一个完整的请求，处理及响应事件。 | 0     |
| lang    | string | 语音流的语言，目前支持zh-CN，en-US。   | zh-CN |
| cdomain | string | 设备当前应用对应的domain           | 空     |

### NlpResponse

| 参数   | 类型     | 描述       |
| ---- | ------ | -------- |
| nlp  | string | nlp识别的结果 |

### TtsRequest

| 参数     | 类型        | 描述                   | 默认值  |
| ------ | --------- | -------------------- | ---- |
| header | TtsHeader | 配置如何将text转换成voice语音流 | 无    |
| text   | string    | 需要转换的text文本          | 无    |


### TtsHeader
PCM格式为 RIFF (little-endian) data, WAVE audio, Microsoft PCM, 16 bit, mono 24000 Hz语音。
注意，这里的PCM与Asr的PCM格式不一致，因此不能将tts的输出直接作为asr的输入。

| 参数        | 类型     | 描述                          | 默认值  |
| --------- | ------ | --------------------------- | ---- |
| id        | int32  | 唯一标识，用于跟踪一个完整的请求，处理及响应事件。   | 0    |
| declaimer | string | 发音者，如"zh","zhangsan","rose" | "zh" |
| codec     | string | 语音流的编码，目前支持PCM，OPU，OPU2。    | PCM  |

### TtsResponse

| 参数    | 类型     | 描述            |
| ----- | ------ | ------------- |
| text  | string | voice语音中包含的文字 |
| voice | bytes  | 合成的voice语音    |

### VoiceSpeechRequest

| 参数     | 类型        | 描述                   | 默认值  |
| ------ | --------- | -------------------- | ---- |
| header | TtsHeader | 配置如何将text转换成voice语音流 | 无    |
| voice  | bytes     | 需要识别的voice语音流        | 无    |

### TexeSpeechRequest

| 参数     | 类型        | 描述                   | 默认值  |
| ------ | --------- | -------------------- | ---- |
| header | TtsHeader | 配置如何将text转换成voice语音流 | 无    |
| asr    | string    | 需要识别的asr文本           | 无    |


### SpeechHeader

| 参数      | 类型     | 描述                        | 默认值   |
| ------- | ------ | ------------------------- | ----- |
| id      | int32  | 唯一标识，用于跟踪一个完整的请求，处理及响应事件。 | 0     |
| lang    | string | 语音流的语言，目前支持zh-CN，en-US。   | zh-CN |
| codec   | string | 语音流的编码，目前支持PCM，OPU，OPU2。  | PCM   |
| vt      | string | 激活词，即用于唤醒设备的名字，如"若琪"。     | 空     |
| cdomain | string | 设备当前的domain信息。            | 空     |
| device  | string | 设备上的状态信息，为json结构。         | 空     |

### SpeechResponse

| 参数     | 类型     | 描述             |
| ------ | ------ | -------------- |
| asr    | string | asr实时识别的结果     |
| nlp    | string | nlp识别的结果       |
| action | string | cloud app处理的结果 |

