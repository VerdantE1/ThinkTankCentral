# Service文档

## base.go

### 为用户发送邮箱验证码进行验证/SendEmailVerificationCode

```go
func (baseService *BaseService) SendEmailVerificationCode(c *gin.Context, to string) error
```

#### 目的

该算法的目的是为用户发送邮箱验证码进行验证，确保用户邮箱的真实性和安全性。这种功能通常用在用户注册、修改邮箱、找回密码等场景，目的是通过验证码验证用户对指定邮箱的控制权。

#### 流程

1. **生成验证码和过期时间**：
   - 调用 `utils.GenerateVerificationCode(6)` 生成一个 6 位数字的随机验证码，用于邮箱验证。
   - 设置验证码的有效期为 5 分钟，通过 `time.Now().Add(5 * time.Minute).Unix()` 获取验证码的过期时间。
2. **存储验证码相关信息到会话（Session）**：
   - 获取当前请求的会话对象：`sessions.Default(c)`。
   - 将生成的验证码、目标邮箱地址、验证码的过期时间分别存储在会话中，键名为 `verification_code`、`email` 和 `expire_time`。
   - 调用 `session.Save()` 将会话信息保存到存储介质中。
3. **生成邮件内容**：
   - 定义邮件的主题：`subject = "您的邮箱验证码"`。
   - 构造邮件的正文（`body`）内容 ，包括：
     - 向用户致以问候。
     - 感谢用户注册或使用服务。
     - 提供生成的验证码，提醒用户验证码有效期为 5 分钟。
     - 提供支持团队的联系信息，以便用户在遇到问题时能获得帮助。
4. **发送邮件**：
   - 调用 `utils.Email(to, subject, body)` 函数，将邮件发送到指定的目标邮箱地址 `to`，包括主题和正文内容。
5. **返回结果**：
   - 函数最后返回 `nil`，表示邮箱验证码已成功生成并发送（假设没有错误处理）。

## es_index.go（不懂）

### 在 Elasticsearch 中创建一个新的索引/IndexCreate

```go
func (esService *EsService) IndexCreate(indexName string, mapping *types.TypeMapping) error
```

#### 目的

该方法的目的是在 Elasticsearch 中创建一个新的索引，并为该索引设置映射（Mapping）。映射定义了索引中字段的结构、数据类型及其他规则，确保存储和检索数据时符合预期的格式和逻辑。

#### 流程

1. **接收参数**：
   - `indexName`：指定要创建的 Elasticsearch 索引的名称。
   - `mapping`：索引的映射结构（`*types.TypeMapping`），定义该索引中字段的类型、属性及存储规则。
2. **调用ESClient的 `Indices.Create` 方法**：
   - 使用 `global.ESClient.Indices.Create(indexName)` 调用 ESClient，以创建指定名称的索引。
   - 调用 `.Mappings(mapping)` 设置索引的字段映射规则。
3. **执行创建操作**：
   - 调用 `.Do(context.TODO())` 方法执行创建索引的请求。
   - 如果请求执行成功，则索引会被成功创建。
4. **返回错误信息**：
   - 检查创建操作的执行结果：
     - 如果成功，返回 `nil`，表示创建索引操作成功。
     - 如果失败，返回错误对象 `err`，以便调用方处理异常情况。

### 从 Elasticsearch 中删除指定的索引/IndexDelete

```go
func (esService *EsService) IndexDelete(indexName string) error
```

#### 目的

`IndexDelete` 方法的目的在于从 Elasticsearch 中删除指定的索引（Index）。删除索引可以清理不再需要的数据，释放存储空间，或为重新创建索引做准备。此功能在数据管理和系统维护中非常重要。

#### 流程

1. **接收参数**：
   - `indexName`：指定要删除的 Elasticsearch 索引的名称。
2. **调用 ESClient的 `Indices.Delete` 方法**：
   - 使用 `global.ESClient.Indices.Delete(indexName)` 调用 ESClient，指定要删除的索引名称。
3. **执行删除操作**：
   - 调用 `.Do(context.TODO())` 方法执行删除请求。
   - ESClient 会根据传入的索引名称，尝试删除对应的索引。
4. **返回错误信息**：
   - 检查删除操作的执行结果：
     - 如果删除成功，返回 `nil`，表示索引删除成功。
     - 如果删除失败（例如索引不存在或权限不足），返回错误对象 `err`。

### 检查指定的 Elasticsearch 索引是否存在/IndexExists

```go
func (esService *EsService) IndexExists(indexName string) (bool, error)
```

#### 目的

`IndexExists` 方法的目的是检查指定的 Elasticsearch 索引是否存在。通过此功能，开发者可以判断目标索引是否已被创建，从而避免重复创建或误操作。例如，在动态创建索引时，可以先检查索引是否存在，再决定是否执行创建操作。

#### 流程

1. **接收参数**：
   - `indexName`：指定需要检查的 Elasticsearch 索引名称。
2. **调用 ESClient的 `Indices.Exists` 方法**：
   - 使用 `global.ESClient.Indices.Exists(indexName)` 调用 ESClient，传入要检查的索引名称。
3. **执行检查操作**：
   - 调用 `.Do(context.TODO())` 方法执行检查请求。
   - ESClient 会返回一个布尔值，指示索引是否存在。
4. **返回结果**：
   - 如果索引存在：
     - 返回 `true` 和 `nil`，表示索引存在，且没有错误发生。
   - 如果索引不存在：
     - 返回 `false` 和 `nil`，表示索引不存在，且没有错误发生。
   - 如果发生其他错误（例如网络问题、权限不足等），返回 `false` 和错误对象 `err`。

## gaode.go

### 根据用户提供的 IP 地址获取对应的地理位置信息/GetLocationByIP

```go
func (gaodeService *GaodeService) GetLocationByIP(ip string) (other.IPResponse, error)
```

#### 目的

`GetLocationByIP` 方法的目的是通过调用高德地图的 IP 定位 API，根据用户提供的 IP 地址获取对应的地理位置信息（如城市、区域等）。此功能可以用于位置识别、个性化推荐、地理限制等场景。

#### 流程

1. **初始化数据结构**：
   - 定义返回值 `data` 为 `other.IPResponse` 类型，用于存储高德 API 的响应结果。
2. **构造请求所需的参数**：
   - 从配置中获取高德地图的 API 密钥：`key := global.Config.Gaode.Key`。
   - 设置请求的 URL：`urlStr := "https://restapi.amap.com/v3/ip"`。
   - 定义请求方法为 `GET`。
   - 创建请求参数：`params := map[string]string`，包括：
     - `ip`：用户提供的 IP 地址。
     - `key`：高德地图 API 的密钥，用于授权。
3. **发送 HTTP 请求**：
   - 调用工具函数 `utils.HttpRequest`，传入请求的 URL、方法、参数等，发送 HTTP 请求。
   - 如果请求返回错误，直接返回 `data` 和错误信息。
4. **检查 HTTP 响应状态码**：
   - 如果响应状态码不是 `200 OK`，返回 `data` 和包含错误信息的 `fmt.Errorf`。
5. **读取响应数据**：
   - 使用 `io.ReadAll` 读取 HTTP 响应的 Body。
   - 如果读取失败，直接返回错误。
6. **解析 JSON 数据**：
   - 使用 `json.Unmarshal` 将响应的 JSON 数据解析到 `data`（`other.IPResponse` 类型）。
   - 如果解析失败，返回错误。
7. **返回结果**：
   - 如果一切正常，返回解析后的 `data` 和 `nil`。
   - 如果任何步骤出错，返回错误信息。

### 根据提供的城市编码（adcode）获取该城市的实时天气信息/GetWeatherByAdcode

```go
func (gaodeService *GaodeService) GetWeatherByAdcode(adcode string) (other.Live, error)
```

#### 目的

`GetWeatherByAdcode` 方法的目的是通过调用高德地图的天气查询 API，根据提供的城市编码（`adcode`）获取该城市的实时天气信息。此功能可用于天气查询、智能推荐、出行规划等场景。

#### 流程

1. **初始化数据结构**：
   - 定义返回数据结构 `data` 为 `other.WeatherResponse` 类型，用于存储高德地图 API 的天气查询响应结果。
2. **构造请求参数**：
   - 从配置中获取高德地图 API 的密钥：`key := global.Config.Gaode.Key`。
   - 设置请求的 URL：`urlStr := "https://restapi.amap.com/v3/weather/weatherInfo"`。
   - 定义请求方法为 `GET`。
   - 创建请求的参数 `params := map[string]string`，包括：
     - `city`: 城市编码（`adcode`）。
     - `key`: 高德地图 API 的密钥。
3. **发送 HTTP 请求**：
   - 调用工具函数 `utils.HttpRequest`，传入 URL、方法、参数等，向高德 API 发送 HTTP 请求。
   - 如果请求失败，直接返回空的 `other.Live` 和错误信息。
4. **检查 HTTP 响应状态码**：
   - 如果响应状态码不是 `200 OK`，返回错误信息（包含状态码）。
5. **读取响应数据**：
   - 使用 `io.ReadAll` 读取 HTTP 响应的 Body。
   - 如果读取失败，直接返回错误。
6. **解析 JSON 数据**：
   - 使用 `json.Unmarshal` 将响应的 JSON 数据解析为 `data`（`other.WeatherResponse` 类型）。
   - 如果解析失败，返回错误。
7. **检查返回的天气数据**：
   - 判断 `data.Lives`是否为空：
     - 如果 `data.Lives` 为空，说明 API 没有返回实时天气信息，返回错误。
     - 如果有数据，提取 `data.Lives[0]`，即当天的实时天气数据。
8. **返回结果**：
   - 如果一切正常，返回 `data.Lives[0]` 和 `nil`。
   - 如果任何步骤出错，返回空的 `other.Live` 和错误信息。

## jwt.go

用户登录（发送用户名和密码）
           │
          ▼
    服务器验证凭据
          │
          ├── 验证失败：返回错误
          │
          └── 验证成功：生成 JWT 并返回
                       │
                      ▼
       客户端存储 JWT（Local Storage/Cookies）
                      │
                     ▼
      客户端携带 JWT 访问受保护资源
                     │
                    ▼
      服务器验证 JWT 签名和有效性
                  │
                  ├── 验证失败：返回 401 未授权
                  │
                  └── 验证成功：允许访问资源

### 将 JWT存储到 Redis 中，并设置一个过期时间/SetRedisJWT

```go
func (jwtService *JwtService) SetRedisJWT(jwt string, uuid uuid.UUID) error
```

#### 目的

`SetRedisJWT` 方法的目的是将 JWT（JSON Web Token）存储到 Redis 中，并设置一个过期时间。通过此方法，可以实现对用户登录状态的管理，例如让 JWT 在一定时间后自动失效，从而提高系统的安全性和可维护性。

#### 流程

1. **解析配置中的 JWT 过期时间**：
   - 从全局配置中读取 `Jwt.RefreshTokenExpiryTime`，这是一个表示 JWT 刷新令牌过期时间的字符串。
   - 调用工具函数 `utils.ParseDuration` 将配置的字符串格式转换为 `time.Duration`，以便用于 Redis 设置过期时间。
   - 如果解析失败，返回解析错误。
2. **将 JWT 存储到 Redis 中**：
   - 使用 `uuid.String()` 作为 Redis 的键，将 JWT 作为值存储到 Redis 中。
   - 调用 Redis 的 Set 方法，并传入：
     - 键：`uuid.String()`。
     - 值：`jwt`（待存储的 JWT）。
     - 过期时间：`dr`（解析得到的过期时间）。
3. **检查存储结果**：
   - 调用 Redis 操作的 `Err()` 方法，检查存储操作是否成功。
   - 如果 Redis 操作成功，返回 `nil`。
   - 如果 Redis 操作失败，返回错误对象。

### 通过指定的 UUID 从 Redis 中获取对应的 JWT/GetRedisJWT

```go
func (jwtService *JwtService) GetRedisJWT(uuid uuid.UUID) (string, error)
```

#### 目的

`GetRedisJWT` 方法的目的是通过指定的 `UUID` 从 Redis 中获取对应的 JWT（JSON Web Token）。该方法通常用于验证用户的登录状态或从 Redis 中取出已经存储的 JWT，以便进一步操作（如验证或刷新令牌）。

#### 流程

1. **接收输入**：
   - 方法接收一个 `uuid.UUID` 类型的参数，表示需要查找的 Redis 键名。
2. **将 UUID 转换为字符串**：
   - 调用 `uuid.String()` 方法，将 `UUID` 转换为字符串格式，因为 Redis 的键名是字符串类型。
3. **从 Redis 获取值**：
   - 调用 `global.Redis.Get(uuid.String()).Result()` 方法：
     - `uuid.String()` 是 Redis 的键名。
     - `global.Redis.Get()` 方法用于从 Redis 中获取指定键的值。
     - `.Result()`方法返回两个值：
       - 获取到的值（JWT 字符串）。
       - 错误对象（如果获取失败）。
4. **返回结果**：
   - 如果 Redis 成功返回值，则返回对应的 `JWT` 和 `nil`。
   - 如果 Redis 返回错误（如键不存在或 Redis 操作失败），则返回空字符串和错误信息。

### 将指定的 JWT 添加到黑名单/JoinInBlacklist

```go
func (jwtService *JwtService) JoinInBlacklist(jwtList database.JwtBlacklist) error
```

#### 目的

`JoinInBlacklist` 方法的主要目的是在 **JWT 黑名单机制** 中，将指定的 JWT 添加到黑名单。通过将 JWT 标记为无效，这种机制可以实现令牌的主动失效（例如用户登出、账户被禁用等场景），从而提升系统的安全性。

具体而言：

1. 将指定的 JWT 记录插入到数据库中的黑名单表，确保黑名单信息持久化。
2. 将 JWT 加入内存中的黑名单缓存，以便快速验证 JWT 是否已被拉黑，提升系统性能。

#### 流程

1. **接收输入**：
   - 方法接收一个 `JwtBlacklist` 类型的参数（通常包含 JWT 和相关信息），表示需要加入黑名单的 JWT。
2. **将 JWT 插入数据库**：
   - 调用 `global.DB.Create(&jwtList)` 方法，将黑名单记录插入到数据库的黑名单表中。
   - 如果插入失败，返回错误信息。
3. **将 JWT 添加到内存缓存**：
   - 调用 `global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})`方法，将 JWT 添加到内存中的黑名单缓存。
     - 键：`jwtList.Jwt`（JWT 字符串）。
     - 值：`struct{}{}`（空结构体），表示该 JWT 已被拉黑。
   - 使用内存缓存可以加速黑名单的查询，减少对数据库的依赖。
4. **返回结果**：
   - 如果数据库插入和缓存更新都成功，返回 `nil`，表示操作成功。
   - 如果数据库插入失败，返回对应的错误信息。

### 检查某个 JWT 是否已经被添加到黑名单中/IsInBlacklist

```go
func (jwtService *JwtService) IsInBlacklist(jwt string) bool
```

#### 目的

`IsInBlacklist` 方法的目的是检查某个 JWT 是否已经被添加到黑名单中。通过这种机制，可以在请求验证阶段快速判断某个 JWT 是否失效，防止被拉黑的令牌继续使用。

- 主要用途：用于实现主动令牌失效机制。
  - 如果 JWT 在黑名单中，则表明它已失效（例如用户已登出、管理员禁用账户等）。
  - 如果 JWT 不在黑名单中，则表明它是有效的。

该方法通过查询内存缓存（`global.BlackCache`）快速判断 JWT 是否被拉黑，避免频繁访问数据库，提高性能。

#### 流程

1. **接收输入**：
   - 方法接收一个 `jwt`（字符串类型），表示需要检查的 JSON Web Token。
2. **查询内存缓存**：
   - 调用 `global.BlackCache.Get(jwt)` 方法，在黑名单缓存中查找该 JWT。
   - `global.BlackCache` 是一个缓存实例，通常使用内存缓存（如 Go 的 `sync.Map` 或第三方库 `cache`），存储已被拉黑的 JWT。
3. **判断是否存在**：
   - 如果缓存中存在该 JWT：
     - 返回 `true`，表示该 JWT 在黑名单中。
   - 如果缓存中不存在该 JWT：
     - 返回 `false`，表示该 JWT 不在黑名单中。
4. **返回结果**：
   - 返回布尔值 `true` 或 `false`，分别表示令牌是否在黑名单中。

### 从数据库中加载所有已经存储的 JWT 黑名单到缓存中/LoadAll

```go
func LoadAll()
```

#### 目的

`LoadAll` 方法的主要目的是从数据库中加载所有已经存储的 JWT 黑名单，并将它们添加到内存缓存中（`global.BlackCache`）。通过将黑名单存储在内存中，可以在请求验证阶段快速判断某个 JWT 是否已失效，从而提高系统性能，减少对数据库的访问。

#### 流程

1. **初始化变量**：
   - 定义一个字符串切片 `data`，用于存储从数据库中查询到的黑名单 JWT。
2. **从数据库加载黑名单**：
   - 调用 `global.DB.Model(&database.JwtBlacklist{}).Pluck("jwt", &data)`：
     - 查询 `JwtBlacklist` 表中所有的 `jwt` 字段值，并存储到 `data` 切片中。
     - 如果查询失败（如数据库连接失败或表数据异常），记录错误日志并终止操作。
3. **将黑名单加载到缓存**：
   - 遍历 `data` 切片，将每个 JWT 添加到内存缓存 `global.BlackCache` 中。
   - 调用 `SetDefault(jwt, struct{}{})` 方法：
     - 使用 `jwt` 作为键。
     - 使用空结构体 `struct{}{}` 作为值（占用最小内存）。
4. **完成加载**：
   - 当所有 JWT 都成功添加到缓存中后，结束操作。

## qq.go

### 通过 Authorization Code获取用户的 Access Token/GetAccessTokenByCode

```go
func (qqService *QQService) GetAccessTokenByCode(code string) (other.AccessTokenResponse, error)
```

#### **目的**

`GetAccessTokenByCode` 方法的目的是通过 **Authorization Code** 从 QQ 的 OAuth 2.0 接口（https://graph.qq.com/oauth2.0/token）获取用户的 **Access Token**。Access Token 是 QQ API 的核心凭据，用于授权后续的 API 调用，例如获取用户信息、获取 OpenID 等。

#### **流程**

**1. 准备工作**

- 输入：
  - 方法接收一个 `code` 参数，这是通过 QQ OAuth 登录时返回的授权码（Authorization Code）。它是临时的，并且会过期。
- 配置：
  - 从 `global.Config.QQ`中读取 QQ 应用的相关配置信息，包括：
    - `AppID`：QQ 应用的唯一标识。
    - `AppKey`：QQ 应用的密钥，用于认证请求。
    - `RedirectURI`：QQ 应用的回调地址，和 QQ 平台配置的回调地址必须一致。

**2. 构造请求**

- URL：
  - 将请求地址设置为 `https://graph.qq.com/oauth2.0/token`。
- HTTP 方法：
  - 使用 `GET` 方法发起请求。
- 请求参数：
  - 构造请求参数 `params`，包含以下字段：
    - `grant_type`: 固定为 `"authorization_code"`，表示通过 Authorization Code 获取 Access Token。
    - `client_id`: 应用的 `AppID`。
    - `client_secret`: 应用的 `AppKey`。
    - `code`: 授权码（Authorization Code）。
    - `redirect_uri`: 回调地址（必须和 QQ 平台配置的一致）。
    - `fmt`: 设置为 `"json"`，表示返回 JSON 格式响应。
    - `need_openid`: 设置为 `"1"`，表示需要返回 `openid`。

**3. 发起请求**

- 调用工具方法 `utils.HttpRequest`发起 HTTP 请求：
  - **URL**：`https://graph.qq.com/oauth2.0/token`。
  - **HTTP 方法**：`GET`。
  - **请求参数**：通过 `params` 传递。
  - **请求头**：为简化，未设置额外的 HTTP 头。

**4. 检查响应**

- 检查错误：
  - 如果请求失败（`HttpRequest` 方法返回错误），直接返回错误信息。
- 检查状态码：
  - 如果响应的 HTTP 状态码不是 `200 OK`，返回错误信息，包含具体的状态码。

**5. 解析响应**

- 读取响应体：
  - 调用 `io.ReadAll(res.Body)` 读取响应体的字节数据。
- JSON 解析：
  - 调用 `json.Unmarshal` 将 JSON 格式的响应体解析为 `other.AccessTokenResponse` 结构体。
  - 如果解析失败，返回错误信息。

**6. 返回结果**

- 返回解析后的 `AccessTokenResponse` 数据和 `nil` 错误，表示请求成功。

### 通过 **Access Token** 和 **OpenID** 获取登录用户的详细信息/GetUserInfoByAccessTokenAndOpenid

```go
func (qqService *QQService) GetUserInfoByAccessTokenAndOpenid(accessToken, openID string) (other.UserInfoResponse, error)
```

#### **目的**

`GetUserInfoByAccessTokenAndOpenid` 方法的主要目的是通过 **Access Token** 和 **OpenID** 调用 QQ 的用户信息接口（https://graph.qq.com/user/get_user_info），获取登录用户的详细信息（例如昵称、头像、性别等）。该方法是实现 QQ 登录功能的核心部分之一，通过 QQ OAuth 授权后，使用此接口获取用户的相关信息。

#### 流程

**1. 准备工作**

- 输入参数：
  - `accessToken`：授权服务器返回的 Access Token，用于访问受保护资源。
  - `openID`：用户在 QQ 平台上的唯一标识符，用于标记具体用户。
- 配置：
  - 从全局配置 `global.Config.QQ` 中获取应用的 `AppID`，用于标识调用接口的应用。

**2. 构造请求**

- URL：
  - 设置请求地址为 `https://graph.qq.com/user/get_user_info`。
- HTTP 方法：
  - 使用 `GET` 方法发起请求。
- 请求参数：
  - 构造请求参数 `params`，包含以下字段：
    - `access_token`: 使用 OAuth 授权获取的 Access Token。
    - `oauth_consumer_key`: 应用的 `AppID`，表示调用接口的客户端。
    - `openid`: 用户的 OpenID，唯一标识用户。
  - 这些参数是接口调用的必需内容，缺失或错误会导致请求失败。

**3. 发起请求**

- 调用工具方法 `utils.HttpRequest` 发起 HTTP 请求：
  - **URL**：`https://graph.qq.com/user/get_user_info`。
  - **HTTP 方法**：`GET`。
  - **请求参数**：通过 `params` 传递。
  - **请求头**：未设置额外的 HTTP 头。

**4. 检查响应**

- 检查错误：
  - 如果请求失败（`HttpRequest` 方法返回错误），直接返回错误信息。
- 检查状态码：
  - 如果响应的 HTTP 状态码不是 `200 OK`，返回错误信息并包含具体的状态码。

**5. 解析响应**

- 读取响应体：
  - 调用 `io.ReadAll(res.Body)` 读取响应体的字节数据。
- JSON 解析：
  - 调用 `json.Unmarshal` 将 JSON 格式的响应体解析为 `other.UserInfoResponse` 结构体。
  - 如果解析失败，返回错误信息。

**6. 返回结果**

- 返回解析后的 `UserInfoResponse` 数据和 `nil` 错误，表示请求成功。

## user.go

### 用户注册/Register

```go
func (userService *UserService) Register(u database.User) (database.User, error)
```

#### 目的

`Register` 方法的主要目的是为用户提供注册功能。通过该方法，用户可以使用电子邮件注册账户，生成必要的用户信息（如密码加密、分配默认头像等），并将用户数据写入数据库。

该方法确保了用户注册的合法性，例如检查电子邮件是否已经被注册，避免重复注册。

#### 流程

**1. 检查电子邮件是否已注册**

- 调用 `global.DB.Where("email = ?", u.Email).First(&database.User{})`：
  - 在数据库中查找是否存在与传入的电子邮件地址相同的用户记录。
  - 如果找到记录（即电子邮件已被注册），返回错误提示，告知用户电子邮件已注册。

**2. 处理用户数据**

- 对用户提供的数据（如密码）进行处理，以确保安全性和规范性：
  - 密码加密：
    - 调用 `utils.BcryptHash(u.Password)` 将用户提供的明文密码加密为哈希值，确保密码不会以明文存储在数据库中。
  - 生成 UUID：
    - 调用 `uuid.NewV4()` 为用户生成一个全局唯一标识符（UUID），用于标识用户。
  - 分配默认头像：
    - 设置用户的头像路径为 `/image/avatar.jpg`，提供一个默认头像。
  - 设置用户角色：
    - 将用户角色设置为普通用户（`appTypes.User`），表示该用户是普通注册用户。
  - 设置注册来源：
    - 将注册来源标记为电子邮件注册（`appTypes.Email`），以便后续统计或分析。

**3. 将用户数据保存到数据库**

- 调用 `global.DB.Create(&u)`将用户数据插入到数据库。
  - 如果插入失败（例如数据库发生错误），返回错误信息。
  - 如果插入成功，返回完整的用户记录。

**4. 返回结果**

- 如果所有步骤成功执行，返回已注册的用户对象和 `nil` 错误，表示注册成功。
- 如果任何步骤失败，返回空用户对象和具体的错误信息。

### 用户邮件登录/EmailLogin

```go
func (userService *UserService) EmailLogin(u database.User) (database.User, error)
```

#### 目的

`EmailLogin` 方法的目的是实现基于电子邮件和密码的用户登录功能。通过验证用户输入的电子邮件和密码是否正确，返回对应的用户数据或错误信息。此方法主要用于用户登录系统，通过数据库查询和密码校验，确保只有正确的用户凭证能够通过验证并登录成功。

#### 流程

**1. 根据电子邮件查找用户**

- 操作：
  - 调用 `global.DB.Where("email = ?", u.Email).First(&user)`，在数据库中查询是否存在与用户输入的电子邮件匹配的记录。
- 结果：
  - 如果查询成功，说明该电子邮件对应的用户存在，继续下一步校验密码。
  - 如果查询失败，直接返回查询错误（如用户不存在）。

**2. 校验用户密码**

- 操作：
  - 调用 `utils.BcryptCheck(u.Password, user.Password)`，验证用户输入的密码（`u.Password`）与数据库中存储的加密密码（`user.Password`）是否匹配。
- 结果：
  - 如果校验成功（`ok == true`），说明密码正确，登录验证通过。
  - 如果校验失败（`ok == false`），返回错误提示 `"incorrect email or password"`，告知用户登录失败。

**3. 返回结果**

- 如果电子邮件存在且密码校验成功：
  - 返回对应的用户数据（`user`）和 `nil` 错误，表示登录成功。
- 如果电子邮件不存在或密码校验失败：
  - 返回空用户对象（`database.User{}`）和相应错误信息。

### 通过QQ登录/QQLogin

#### 目的

`QQLogin` 方法的主要目的是通过 QQ 登录功能，让用户使用 QQ 的 OAuth 认证登录系统。它通过 QQ 平台获取用户的基本信息，并在本地系统中创建或更新用户记录。此方法包括两种场景：

1. 如果用户已经存在于数据库中，直接返回用户信息。
2. 如果用户不存在，则调用 QQ API 获取用户QQ信息并以此创建新用户。

#### 流程

**1. 从数据库中查找用户**

- 操作：
  - 调用 global.DB.Where("openid = ?", accessTokenResponse.Openid).First(&user)：
    - 根据 QQ 返回的 OpenID 查询用户是否已经存在于数据库中。
  - 结果：
    - 如果用户记录存在，直接返回用户信息。
    - 如果发生其他数据库错误（非用户不存在），返回错误信息。
    - 如果用户不存在（返回 `gorm.ErrRecordNotFound`），进入下一步创建新用户。

**2. 调用 QQ 用户信息接口**

- 操作：

  - 调用 QQService.GetUserInfoByAccessTokenAndOpenid(accessTokenResponse.AccessToken, accessTokenResponse.Openid)：

    - 使用 Access Token 和 OpenID 调用 QQ 的用户信息接口，获取用户的详细信息（如昵称、头像）。

  - 结果：

    - 如果接口调用失败，返回错误信息。

    - 如果调用成功，获取 QQ 用户的详细信息（如昵称和头像），用于创建新用户记录。

**3. 创建新用户**

- 操作：
  - 如果用户不存在于数据库中，初始化用户数据并保存到数据库：
    - **生成 UUID**：为用户生成全局唯一标识符（UUID）。
    - **设置用户名**：使用 QQ 用户的昵称作为用户名。
    - **设置 OpenID**：保存 QQ 返回的 OpenID，用于标识用户。
    - **设置头像**：使用 QQ 用户的头像 URL。
    - **设置用户角色**：将用户角色设置为普通用户（`appTypes.User`）。
    - **设置注册方式**：标记用户通过 QQ 注册（`appTypes.QQ`）。
  - 调用 `global.DB.Create(&user)` 将用户数据保存到数据库。
- 结果：
  - 如果保存成功，返回新创建的用户信息。
  - 如果保存失败，返回错误信息。

**4. 返回用户信息**

- 操作：
  - 如果用户已存在或成功创建新用户，返回用户信息。
  - 如果任何步骤失败，返回空用户对象和错误信息。

### 忘记密码后重置密码/ForgotPassword

#### **目的**

`ForgotPassword` 方法的主要目的是为用户提供重置密码功能。当用户忘记密码时，通过此方法可以更新用户的密码。该方法根据用户提供的邮箱地址查找用户记录，并将新密码加密后保存到数据库中。

#### 流程

**1. 根据邮箱查找用户**

- 操作：
  - 调用global.DB.Where("email = ?", req.Email).First(&user)：
    - 在数据库中查询是否存在与请求中提供的电子邮件地址（`req.Email`）匹配的用户记录。
- 结果：
  - 如果查询成功，找到对应的用户记录，进入下一步更新密码。
  - 如果查询失败（例如用户不存在），直接返回错误信息。

**2. 加密新密码**

- 操作：
  - 调用 utils.BcryptHash(req.NewPassword)：
    - 使用加密算法（如 bcrypt）对用户提供的新密码进行哈希加密，避免明文密码存储到数据库中。
- 结果：
  - 将加密后的密码赋值给用户对象的 `Password` 字段。

**3. 更新密码到数据库**

- 操作：
  - 调用 global.DB.Save(&user)：
    - 将更新后的用户记录保存到数据库中，覆盖旧的密码。
- 结果：
  - 如果保存成功，返回 `nil`，表示密码重置成功。
  - 如果保存失败（例如数据库错误），返回错误信息。

UserCard

### 根据用户的 UUID 查询用户资料卡片/UserCard

#### 目的

`UserCard` 方法的主要目的是根据用户的 UUID 查询用户的关键信息，并返回一个精简的用户资料卡片。这个方法主要用于展示用户的公共信息（如用户名、头像、地址、个性签名等），而不暴露敏感信息（如密码或电子邮件）。

#### 流程

**1. 根据 UUID 查询用户**

- 操作：
  - 调用 global.DB.Where("uuid = ?", req.UUID).Select(...).First(&user)：
    - 在数据库中根据传入的 `UUID` 查询用户记录。
    - 使用 `Select` 明确指定查询的字段（`uuid`、`username`、`avatar`、`address`、`signature`），只获取必要的数据。
- 结果：
  - 如果查询成功，则获取用户的精简信息，进入下一步。
  - 如果查询失败（例如 UUID 不存在），返回错误。

**2. 格式化返回数据**

- 操作：
  - 构造 response.UserCard 对象，将查询到的用户信息映射到响应结构体中：
    - 将用户的 `UUID`、`Username`、`Avatar`、`Address` 和 `Signature` 填充到响应对象中。
- 结果：
  - 返回构造好的 `UserCard` 响应对象，供前端或调用方使用。

**3. 返回结果**

- 如果查询成功，返回用户的精简信息（`response.UserCard`）和 `nil` 错误。
- 如果查询失败，返回空的 `response.UserCard` 对象和数据库错误。

### 用户登出功能/Logout

#### 目的

`Logout` 方法的主要目的是实现用户登出功能。通过清除用户的刷新令牌、从 Redis 中移除用户的会话数据以及将 JWT 加入黑名单，确保用户安全退出系统并防止后续的令牌滥用。

#### 流程

**1. 获取用户 UUID**

- 操作：
  - 调用 `utils.GetUUID(c)` 从上下文中提取用户的 UUID。
  - UUID 是用户的唯一标识，用于标记当前用户的会话数据。
- 结果：
  - 成功提取用户 UUID，供后续操作使用。

**2. 获取刷新令牌**

- 操作：
  - 调用 `utils.GetRefreshToken(c)` 从上下文中提取用户的刷新令牌（JWT）。
  - 刷新令牌用于生成新的访问令牌，登出时需要清除。

**3. 清除刷新令牌**

- 操作：
  - 调用 `utils.ClearRefreshToken(c)` 清除用户的刷新令牌（通常是从 Cookie 或 Header 中移除）。
  - 目的是确保客户端不会再持有有效的刷新令牌。

**4. 从 Redis 中删除会话数据**

- 操作：
  - 调用 `global.Redis.Del(uuid.String())` 从 Redis 中删除以用户 UUID 为键的会话信息。
  - Redis 通常用于存储与用户会话相关的数据（如在线状态、临时缓存等）。
- 结果：
  - 如果删除成功，用户在 Redis 中的会话数据将被清除。

**5. 将 JWT 加入黑名单**

- 操作：
  - 调用 `ServiceGroupApp.JwtService.JoinInBlacklist()`，将当前 JWT 刷新令牌加入黑名单。
  - 传入参数为 `database.JwtBlacklist{Jwt: jwtStr}`，将用户的刷新令牌记录在黑名单中。
- 结果：
  - 确保该令牌在其有效期内无法再次使用，防止滥用。

### 重置密码/UserResetPassword

#### 目的

`UserResetPassword` 方法的主要目的是提供用户重置密码的功能。通过验证用户的 **原始密码** 是否正确，然后将用户的密码更新为新的加密密码并保存到数据库中。该方法的目标是确保用户能够安全地更改自己的密码，同时防止未经授权的密码修改操作。

#### 流程

**1. 根据用户 ID 查找用户**

- 操作：
  - 调用 global.DB.Take(&user, req.UserID) 在数据库中查找用户。
    - `req.UserID` 是请求中传递的用户 ID。
- 结果：
  - 如果查询成功，则获取用户记录，进入下一步。
  - 如果查询失败（例如用户不存在），返回错误，终止流程。

**2. 校验原始密码**

- 操作：

  - 调用 utils.BcryptCheck(req.Password, user.Password)校验用户输入的原始密码 req.Password 是否与数据库中存储的加密密码 user.Password

     匹配。

    - `utils.BcryptCheck` 是用于验证 bcrypt 加密密码的工具方法。

- 结果：

  - 如果校验成功（即密码匹配），进入下一步。
  - 如果校验失败，返回错误提示 `"original password does not match the current account"`。

**3. 加密新密码**

- 操作：
  - 调用 `utils.BcryptHash(req.NewPassword)` 对用户提供的新密码 `req.NewPassword` 进行加密。
  - 将加密后的密码赋值给用户对象的 `Password` 字段。

**4. 更新密码到数据库**

- 操作：
  - 调用 `global.DB.Save(&user)` 将更新后的用户记录保存到数据库。
  - 如果保存失败，返回数据库错误。
- 结果：
  - 如果保存成功，方法执行完成，返回 `nil` 表示操作成功。

### 返回用户详细信息/UserInfo

#### **目的**

`UserInfo` 方法的目的是根据用户的唯一标识 (`userID`) 从数据库中查询到用户的详细信息，并返回给调用方。 该方法主要用于获取用户的完整信息，适用于需要展示用户资料或处理用户数据的场景。

#### 流程

**1. 根据用户 ID 查询用户数据**

- 操作：
  - 调用 global.DB.Take(&user, userID)：
    - 在数据库中根据用户的主键 `userID` 查询用户记录。
    - `Take` 方法用于精确查找主键对应的单条记录。
- 结果：
  - 如果查询成功，用户记录被赋值到变量 `user`，进入下一步。
  - 如果查询失败（例如用户不存在或数据库错误），返回错误信息并终止流程。

**2. 返回用户数据**

- 操作：
  - 如果查询成功，返回查询到的用户对象 `user` 和 `nil`（无错误）。
  - 如果查询失败，返回一个空的用户对象 `database.User{}` 和错误信息。
- 结果：
  - 成功时：返回用户数据，供调用方使用。
  - 失败时：返回错误，供调用方处理。

### 用户修改个人信息/UserChangeInfo

#### **目的**

`UserChangeInfo` 方法的目的是为用户提供修改个人信息的功能。通过此方法，系统可以根据用户的请求数据更新用户的资料，例如用户名、头像、地址等字段。

#### 流程

**1. 根据用户 ID 查询用户**

- 操作：
  - 调用 global.DB.Take(&user, req.UserID)：
    - 根据请求中的 `UserID` 在数据库中查找用户记录。
- 结果：
  - 如果查询成功，说明用户存在，进入下一步更新数据。
  - 如果查询失败（例如用户不存在），返回错误信息并终止流程。

**2. 更新用户信息**

- 操作：
  - 调用 global.DB.Model(&user).Updates(req)：
    - 使用 GORM 的 `Updates` 方法，根据 `req` 中的字段值更新用户记录。
    - 只更新 `req` 结构体中提供的字段，其他字段保持不变。
- 结果：
  - 如果更新成功，返回 `nil`，表示操作成功。
  - 如果更新失败（例如数据库错误），返回错误信息。

### 根据用户的 IP 地址获取当前天气信息/UserWeather

#### 目的

`UserWeather` 方法的目的是根据用户的 IP 地址获取当前天气信息。为了提高效率和减少对外部 API 的调用次数：

1. **优先从 Redis 缓存中读取天气数据**。
2. 如果缓存中没有数据，则调用高德 API 获取天气信息，并将结果存入 Redis 缓存中。

该方法的目标是高效、动态地为用户提供准确的天气信息。

### 流程

**1. 从 Redis 中获取天气数据**

- 操作：
  - 调用 global.Redis.Get("weather-" + ip).Result()：
    - 根据用户的 IP 地址拼接出 Redis 的键名（`weather-<ip>`）。
    - 尝试从 Redis 缓存中获取对应的天气信息。
- 结果：
  - 如果数据存在，直接返回缓存中的天气信息。
  - 如果数据不存在（或 Redis 查询失败），进入下一步调用高德 API 获取天气数据。

**2. 调用高德 API 获取天气数据**

- 操作：
  - 获取地理位置：
    - 调用 ServiceGroupApp.GaodeService.GetLocationByIP(ip)：
      - 使用用户的 IP 地址调用高德 IP 定位 API，获取用户所在的行政区域编码（`Adcode`）。
      - 如果调用失败，返回错误信息。
  - 获取天气信息：
    - 调用 ServiceGroupApp.GaodeService.GetWeatherByAdcode(ipResponse.Adcode)：
      - 使用获取到的行政区域编码调用高德天气 API，获取当前区域的天气信息。
      - 如果调用失败，返回错误信息。

**3. 生成天气描述字符串**

- 操作：
  - 根据高德天气 API 返回的数据，生成用户可读的天气描述字符串：
    - 包括省份、城市、天气、温度、风向、风级和湿度等信息。

**4. 将天气数据存入 Redis**

- 操作：
  - 调用 global.Redis.Set("weather-"+ip, weather, time.Hour*1).Err()：
    - 将生成的天气描述字符串存入 Redis，键名为 `weather-<ip>`。
    - 设置缓存过期时间为 1 小时（`time.Hour*1`），避免旧数据长期占用缓存空间。
  - 如果 Redis 存储失败，返回错误信息。

**5. 返回结果**

- 操作：
  - 如果 Redis 缓存命中，直接返回缓存中的天气数据。
  - 如果缓存未命中，返回从高德 API 获取并存储到 Redis 的天气数据。
- 结果：
  - 成功时：返回天气描述字符串。
  - 失败时：返回错误信息。

### 统计一定时间范围内用户的登录次数和注册次数/UserChart

#### 目的

`UserChart` 方法的主要目的是统计一定时间范围内用户的登录次数和注册次数，并将结果以图表所需的数据格式返回。 该方法适用于前端需要可视化用户行为趋势，帮助分析用户的活跃度和注册量变化。

#### 流程

**1. 构建查询条件**

- 操作：

  - 使用 GORM 的 Where 方法构建查询条件，筛选出最近 req.Date 天的数据：

    ```go
    where := global.DB.Where(fmt.Sprintf("date_sub(curdate(), interval %d day) <= created_at", req.Date))
    ```

    - `req.Date` 是请求中指定的天数，用于动态设置时间范围。
    - 查询条件确保只筛选 `created_at` 在最近 `req.Date` 天内的数据。

- 结果：

  - 构造出用于后续查询的筛选条件。

**2. 初始化响应结构与日期列表**

- 操作：

  - 初始化响应结构 `res`，用于存放最终结果。

  - 根据指定的日期范围生成日期列表：

    ```go
    startDate := time.Now().AddDate(0, 0, -req.Date)
    for i := 1; i <= req.Date; i++ {
        res.DateList = append(res.DateList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
    }
    ```

    - `startDate` 是时间范围的起点，通过 `AddDate` 方法逐日递增，生成日期。
    - 日期格式化为 `YYYY-MM-DD`，方便前端图表直接使用。

- 结果：

  - 生成的日期列表存储在 `res.DateList` 中，供后续使用。

**3. 获取登录数据**

- 操作：

  - 调用工具方法 utils.FetchDateCounts，统计指定时间范围内的用户登录次数：

    ```go
    loginCounts := utils.FetchDateCounts(global.DB.Model(&database.Login{}), where)
    ```

    - `FetchDateCounts` 方法接收查询构造器和查询条件，返回一个日期与登录次数的映射 (`map[string]int`)。

- 结果：

  - 获取到每一天的登录统计数据。

**4. 获取注册数据**

- 操作：

  - 调用工具方法 utils.FetchDateCounts，统计指定时间范围内的用户注册次数：

    ```go
    registerCounts := utils.FetchDateCounts(global.DB.Model(&database.User{}), where)
    ```

    - 返回一个日期与注册次数的映射 (`map[string]int`)。

- 结果：

  - 获取到每一天的注册统计数据。

**5. 组装响应数据**

- 操作：

  - 遍历日期列表 res.DateList，从登录和注册数据映射中提取对应的统计值，并填充到响应结构中：

    ```go
    for _, date := range res.DateList {
        loginCount := loginCounts[date]
        registerCount := registerCounts[date]
        res.LoginData = append(res.LoginData, loginCount)
        res.RegisterData = append(res.RegisterData, registerCount)
    }
    ```

    - 如果某日期没有对应的统计数据，默认值为 `0`。

- 结果：

  - 将每一天的登录和注册数据分别存入 `res.LoginData` 和 `res.RegisterData`。

**6. 返回结果**

- 操作：
  - 返回组装好的 `res` 数据，包含日期列表、登录数据和注册数据。
- 结果：
  - 成功返回 `response.UserChart` 对象。
  - 如果过程中出现错误（如数据库查询失败），则返回错误信息。

### 实现用户列表的分页查询功能/UserList

#### 目的

`UserList` 方法的主要目的是实现用户列表的分页查询功能。通过动态构建查询条件（如 UUID 和注册状态），并结合分页参数，返回符合条件的用户列表数据以及总记录数。该方法适用于后台管理系统或其他需要用户列表展示和查询的场景。

#### 流程

**1. 初始化查询构造器**

- 操作：
  - 使用全局数据库实例初始化查询构造器：db := global.DB
  - `db` 是 GORM 的查询链，用于构建动态查询条件。

**2. 动态构建查询条件**

- 操作：

  - 根据请求参数 info中的字段动态添加查询条件：

    - 如果 info.UUID不为空，添加 UUID 查询条件：

      ```go
      if info.UUID != nil {
          db = db.Where("uuid = ?", info.UUID)
      }
      ```

      - 用于精确查询指定 UUID 的用户。

    - 如果 info.Register不为空，添加注册状态查询条件：

      ```
      if info.Register != nil {
          db = db.Where("register = ?", info.Register)
      }
      ```

      - 用于筛选注册状态符合条件的用户（如已注册或未注册）。

- 结果：

  - 动态构建出包含所有查询条件的 `db` 对象，用于后续查询。

**3. 配置分页参数**

- 操作：

  - 创建 MySQLOption对象，将分页信息和查询条件一起传递：

    ```go
    option := other.MySQLOption{
        PageInfo: info.PageInfo,
        Where:    db,
    }
    ```

    - `PageInfo` 包含分页参数（如当前页码和每页条数）。
    - `Where` 是动态构建的查询条件。

**4. 调用分页工具方法**

- 操作：

  - 调用 utils.MySQLPagination工具方法执行分页查询：

    ```go
    return utils.MySQLPagination(&database.User{}, option)
    ```

    - 参数 `&database.User{}` 指定查询的表（用户表）。
    - 参数 `option` 包含分页信息和查询条件。
    - MySQLPagination 返回结果包括：
      - 用户列表数据（`interface{}` 类型）。
      - 总记录数（`int64` 类型）。
      - 错误信息（`error` 类型）。

- 结果：

  - 方法返回符合条件的用户列表数据、总记录数和错误信息。

### 冻结指定用户的账户/UserFreeze

#### 目的

`UserFreeze` 方法的主要目的是冻结指定用户的账户，同时将该用户的 JWT（JSON Web Token）加入黑名单，确保用户无法继续使用系统的登录状态。
 该方法适用于后台管理员操作或系统自动禁用用户账户的场景，用于限制违规用户的操作权限。

#### 流程

**1. 根据用户 ID 查找用户并更新冻结状态**

- 操作：
  - 调用 global.DB.Take(&user, req.ID).Update("freeze", true)：
    - 根据传入的用户 ID (`req.ID`)，查询数据库中的用户记录。
    - 更新用户的 `freeze` 字段为 `true`，表示冻结账户。
  - 如果查询或更新操作失败，返回错误。
- 结果：
  - 成功：用户账户的冻结状态被更新。
  - 失败：返回错误信息，终止流程。

**2. 获取用户的 JWT**

- 操作：
  - 调用 ServiceGroupApp.JwtService.GetRedisJWT(user.UUID)：
    - 根据用户的 UUID 从 Redis 中获取用户的 JWT（当前登录令牌）。
    - 如果查询不到令牌（如用户未登录），则返回空字符串。
- 结果：
  - 成功获取用户的 JWT 或空字符串。

**3. 将用户的 JWT 加入黑名单**

- 操作：

  - 如果用户的 JWT 不为空，将其加入黑名单：

    ```go
    _ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
    ```

    - 调用 `JoinInBlacklist` 方法，将该令牌存入黑名单，确保该令牌即使未过期也无法继续使用。

  - 无需返回值或处理错误（因为黑名单操作不影响冻结流程）。

- 结果：

  - 用户的 JWT 被标记为无效。

**4. 返回操作结果**

- 操作：
  - 如果冻结操作和黑名单操作都成功，返回 `nil` 表示操作成功。
  - 如果冻结操作失败，返回错误信息。
- 结果：
  - 返回操作的最终状态。

### 用户解冻/UserUnfreeze

#### 目的

`UserUnfreeze` 方法的目的是解冻指定用户的账户。 通过将用户的 `freeze` 状态设置为 `false`，恢复用户的正常使用权限。 该方法适用于管理员手动解冻被冻结的用户账户或系统自动解除冻结状态的场景。

#### 流程

**1. 根据用户 ID 查询用户**

- 操作：
  - 调用 global.DB.Take(&database.User{}, req.ID)：
    - 根据传入的用户 ID (`req.ID`)，从数据库中查找对应的用户记录。
    - `Take` 方法用于精确查找主键记录。
  - 如果用户记录不存在，则查询失败，返回错误。
- 结果：
  - 成功：找到用户记录，进入下一步。
  - 失败：返回错误信息。

**2. 更新用户冻结状态**

- 操作：

  - 调用 Update("freeze", false)方法，将用户的 freeze 状态更新为 false：

    ```go
    global.DB.Take(&database.User{}, req.ID).Update("freeze", false)
    ```

    - 将用户的冻结状态解冻，表示该用户可以恢复正常使用系统。

  - 如果更新失败（如数据库错误），返回错误信息。

- 结果：

  - 成功：用户冻结状态被更新为解冻。
  - 失败：返回错误信息。

**3. 返回操作结果**

- 操作：
  - 如果查询和更新操作都成功，返回 `nil`，表示解冻成功。
  - 如果查询或更新失败，返回错误信息。
- 结果：
  - 成功：返回 `nil`。
  - 失败：返回错误信息。

### 实现用户登录日志的分页查询功能/UserLoginList

#### 目的

`UserLoginList` 方法的主要目的是实现用户登录日志的分页查询功能。通过动态构建查询条件（如用户 UUID），结合分页参数，返回用户登录日志列表及总记录数。该方法适用于后台管理系统，帮助管理员查看用户登录行为数据。

#### 流程

**1. 初始化查询构造器**

- 操作：
  - 使用全局数据库实例初始化查询构造器：db := global.DB
  - `db` 是 GORM 的查询链，用于后续动态添加查询条件。

**2. 动态构建查询条件**

- 操作：

  - 如果请求参数中包含 UUID，根据 UUID 获取对应的用户 ID 并添加查询条件：

    ```go
    if info.UUID != nil {
        var userID uint
        if err := global.DB.Model(database.User{}).Where("uuid = ?", *info.UUID).Pluck("id", &userID); err != nil {
            return nil, 0, nil // 查询失败返回空结果
        }
        db = db.Where("user_id = ?", userID) // 添加 user_id 的查询条件
    }
    ```

    - 通过 `Pluck` 方法从用户表中获取与 UUID 对应的用户 ID。
    - 将用户 ID 作为条件，筛选登录表中的 `user_id` 字段。

- 结果：

  - 如果 `UUID` 存在且查询成功，添加对应的查询条件。
  - 如果 `UUID` 查询失败，返回空结果。

**3. 配置分页参数**

- 操作：

  - 创建 MySQLOption对象，将分页信息、查询条件和关联字段传递：

    ```go
    option := other.MySQLOption{
        PageInfo: info.PageInfo, // 分页信息
        Where:    db,           // 动态查询条件
        Preload:  []string{"User"}, // 预加载关联的 User 表
    }
    ```

    - `PageInfo` 包含分页所需的参数（如当前页码和每页条数）。
    - `Preload` 指定预加载 `User` 表，用于在返回结果中包含用户信息。

**4. 调用分页工具方法**

- 操作：

  - 调用 utils.MySQLPagination工具方法执行分页查询：

    ```
    return utils.MySQLPagination(&database.Login{}, option)
    ```

    - 参数 `&database.Login{}` 指定查询的表为登录日志表。
    - 参数 `option` 包含分页信息、查询条件和预加载设置。
    - MySQLPagination 返回以下结果：
      - 登录日志列表数据（`interface{}` 类型）。
      - 总记录数（`int64` 类型）。
      - 错误信息（`error` 类型）。

- 结果：

  - 查询成功：返回符合条件的分页数据和总记录数。
  - 查询失败：返回错误信息。