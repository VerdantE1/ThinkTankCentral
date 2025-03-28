# api文档

## base.go

### 生成动态验证码并返回验证码的数据及其对应 ID/Captcha

```go
func (baseApi *BaseApi) Captcha(c *gin.Context)
```

#### 目的

该方法 `Captcha` 的目的是生成一个动态验证码，并返回验证码的图片数据及其对应的唯一 ID。验证码通常用于验证用户的身份或防止恶意请求，例如登录、注册、防止机器人行为等。

#### 流程

1. **创建验证码驱动**：
   - 调用 `base64Captcha.NewDriverDigit` 方法，创建一个数字型验证码的驱动。
   - 驱动的具体配置由 global.Config.Captcha 提供，包括：
     - `Height`：验证码图片的高度。
     - `Width`：验证码图片的宽度。
     - `Length`：验证码的长度（即数字的个数）。
     - `MaxSkew`：验证码字符的最大倾斜度。
     - `DotCount`：验证码图片的干扰点数量。
2. **初始化验证码对象**：
   - 使用 base64Captcha.NewCaptcha(driver, store) 创建验证码对象：
     - `driver` 是上一步生成的验证码驱动。
     - `store` 是一个存储接口，用于保存验证码的 ID 和答案（通常存储在内存中或其他持久化存储中）。
3. **生成验证码**：
   - 调用 `captcha.Generate()` 方法生成验证码。
   - 返回的内容包括：
     - `id`：验证码的唯一标识符，用于后续验证。
     - `b64s`：验证码图片的 Base64 编码字符串，前端可直接显示该图片。
     - 错误信息 `err` 用于判断生成是否成功。
4. **错误处理**：
   - 如果生成验证码失败（`err != nil`），记录错误日志，并通过 `response.FailWithMessage` 返回失败响应给客户端。
5. **返回验证码信息**：
   - 如果生成成功，将 `CaptchaID`（验证码 ID）和 `PicPath`（验证码图片 Base64 数据）封装到结构体 `response.Captcha` 中，并通过 `response.OkWithData` 将结果返回给客户端。

### 接收客户端发送的请求并验证用户输入的动态验证码是否正确/SendEmailVerificationCode

```go
func (baseApi *BaseApi) SendEmailVerificationCode(c *gin.Context)
```

#### 依赖

```go
func (baseService *BaseService) SendEmailVerificationCode(c *gin.Context, to string) error
```

#### 目的

该方法 `SendEmailVerificationCode` 的目的是接收客户端发送的请求，验证用户输入的动态验证码是否正确。如果验证通过，则调用服务层逻辑向指定的邮箱发送验证码。此方法主要用于邮箱验证的场景，比如用户注册、找回密码等。

#### 流程

1. **接收客户端请求并绑定参数**：
   - 定义一个 `request.SendEmailVerificationCode` 结构体（假定该结构体包含 `Email`、`CaptchaID` 和 `Captcha` 字段，用于接收客户端传递的邮箱和验证码相关信息）。
   - 调用 `c.ShouldBindJSON(&req)`，将请求体中的 JSON 数据绑定到 `req` 中。
   - 如果绑定失败（参数不完整或格式错误），返回绑定错误信息并终止流程。
2. **验证用户输入的动态验证码**：
   - 调用 store.Verify(req.CaptchaID, req.Captcha, true) 验证用户输入的验证码是否正确：
     - `req.CaptchaID` 是验证码的唯一标识符。
     - `req.Captcha` 是用户输入的验证码内容。
     - `true` 表示验证时，成功后会删除该验证码（防止重复使用）。
   - 如果验证码验证失败，则直接返回错误信息 "Incorrect verification code"。
3. **调用服务层发送邮箱验证码**：
   - 如果验证码验证成功，调用服务层方法 `baseService.SendEmailVerificationCode`，传递当前上下文 `c` 和用户的邮箱地址 `req.Email`。
   - 如果发送邮件失败（`err != nil`），记录错误日志并返回失败信息 "Failed to send email"。
4. **返回成功响应**：
   - 如果邮箱验证码发送成功，返回成功消息 "Successfully sent email"。

### 向客户端返回一个用于跳转到 QQ 登录授权页面的链接地址/QQLoginURL

```go
func (baseApi *BaseApi) QQLoginURL(c *gin.Context) 
```

#### 目的

`QQLoginURL` 方法的目的是向客户端返回一个用于跳转到 QQ 登录授权页面的链接地址。用户通过该链接可以跳转到 QQ 的 OAuth 授权页面，用于实现第三方登录功能。

#### 流程

1. **获取 QQ 登录地址**：
   - 调用配置项中的方法 `global.Config.QQ.QQLoginURL()`，生成 QQ 登录授权页面的完整 URL。
   - 该 URL 通常包含：
     - 应用的 `client_id`（QQ 开放平台分配的应用 ID）。
     - 回调地址（用户授权后重定向的地址）。
     - 状态参数（用于防止 CSRF 攻击）。
     - 其他必要参数（如授权类型等）。
   - 这些参数通常在项目的配置文件中提前设置好。
2. **返回登录链接**：
   - 使用 `response.OkWithData(url, c)` 将生成的 URL 返回给客户端。
   - 客户端接收到链接后，可以直接跳转到该地址，完成 QQ 登录授权。

## user.go

### 实现用户邮箱注册功能/Register

```go
func (userApi *UserApi) Register(c *gin.Context)
```

#### 目的

`Register` 方法用于实现用户注册功能。通过校验用户提交的注册信息（邮箱、验证码等），验证邮箱验证码的有效性，然后创建用户记录，并在注册成功后生成并返回一个 JWT 令牌，用于后续的身份认证。

#### 流程

1. **接收并绑定请求参数**：
   - 定义一个 `request.Register` 结构体，用于接收客户端提交的注册请求参数（如邮箱、验证码、用户名、密码等）。
   - 调用 `c.ShouldBindJSON(&req)` 将请求体中的 JSON 数据绑定到 `req` 中。
   - 如果参数绑定失败（如参数缺失或格式错误），返回错误信息并终止流程。
2. **校验邮箱地址一致性**：
   - 获取会话中的 `email`（之前发送验证码时存储在会话中）。
   - 判断会话中的 `email` 是否存在，并与用户提交的 `req.Email` 是否一致。
   - 如果不一致，返回错误信息 "This email doesn't match the email to be verified"。
3. **校验验证码有效性**：
   - 获取会话中存储的验证码 `verification_code`。
   - 判断验证码是否存在，并与用户提交的 `req.VerificationCode` 是否一致。
   - 如果验证码不匹配，返回错误信息 "Invalid verification code"。
4. **校验验证码是否过期**：
   - 获取会话中存储的验证码过期时间 `expire_time`。
   - 判断当前时间是否超过验证码过期时间。
   - 如果验证码已过期，返回错误信息 "The verification code has expired, please resend it"。
5. **创建用户记录**：
   - 构造一个 `database.User` 对象，包含用户提交的注册信息（如用户名、密码、邮箱）。
   - 调用服务层方法 `userService.Register` 创建用户记录。
   - 如果创建用户失败，记录错误日志并返回错误信息 "Failed to register user"。
6. **生成 JWT 令牌**：
   - 如果用户注册成功，调用 `userApi.TokenNext(c, user)` 方法为该用户生成一个 JWT 令牌。
   - 令牌生成成功后，将令牌返回给客户端，用于后续的身份认证。

### 用户登录接口入口方法用于根据用户选择的登录方式调用不同的登录方法/Login

```go
func (userApi *UserApi) Login(c *gin.Context)
```

#### 目的

`Login` 方法是一个用户登录接口的入口，用于根据用户选择的登录方式（`flag` 参数）调用不同的登录方法。目前支持两种登录方式：

1. **邮箱登录**：需要校验图形验证码。
2. **QQ 登录**：需要校验授权码。

该方法的目的是为用户提供灵活的多种登录方式，统一管理登录入口，根据 `flag` 参数分发到具体的登录逻辑。

#### 流程

1. **接收并解析登录方式参数**：
   - 从查询参数中获取 `flag` 参数（`c.Query("flag")`），该参数指定用户选择的登录方式。
   - 目前支持的 `flag` 值包括：
     - `email`：邮箱登录。
     - `qq`：QQ登录。
   - 如果未指定 `flag`，默认使用邮箱登录。
2. **根据 `flag` 调用对应的登录方法**：
   - 邮箱登录：
     - 如果 `flag == "email"`，调用 `userApi.EmailLogin(c)` 方法处理邮箱登录逻辑。
     - 邮箱登录通常需要校验用户的图形验证码以及邮箱和密码的正确性。
   - QQ登录：
     - 如果 `flag == "qq"`，调用 `userApi.QQLogin(c)` 方法处理 QQ 登录逻辑。
     - QQ登录通常需要校验用户的 QQ 授权码，验证授权的合法性。
   - 默认处理：
     - 如果 `flag` 参数无效或未提供，则默认调用 `userApi.EmailLogin(c)`，以邮箱登录作为兜底选项。
3. **返回结果**：
   - 每种登录方式的具体逻辑和返回结果由对应的方法 `EmailLogin` 和 `QQLogin` 实现。
   - 如果登录成功，通常会返回用户信息和一个 JWT 令牌。
   - 登录失败时，返回对应的错误信息。

### 邮箱登录/EmailLogin

```go
func (userApi *UserApi) EmailLogin(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) EmailLogin(u database.User) (database.User, error)
```

#### 目的

`EmailLogin` 方法实现了基于邮箱和密码的登录功能。通过校验用户输入的图形验证码和邮箱密码的正确性，确保登录的安全性。如果验证成功，则生成 JWT 令牌并返回给客户端，用于后续的身份认证。

#### 流程

1. **接收并绑定请求参数**：
   - 定义一个 `request.Login` 结构体，用于接收客户端提交的登录请求参数（如邮箱、密码、验证码等）。
   - 调用 `c.ShouldBindJSON(&req)` 将请求体中的 JSON 数据绑定到 `req` 中。
   - 如果参数绑定失败（如缺少字段或格式错误），返回错误信息并终止流程。
2. **校验图形验证码**：
   - 调用 store.Verify(req.CaptchaID, req.Captcha, true) 校验用户提交的图形验证码：
     - `req.CaptchaID` 是验证码的唯一标识符。
     - `req.Captcha` 是用户输入的验证码内容。
     - `true` 表示验证成功后删除该验证码（防止重复使用）。
   - 如果验证码验证失败，返回错误信息 "Incorrect verification code" 并终止流程。
3. **验证邮箱和密码**：
   - 构造一个 `database.User` 对象，包含用户输入的邮箱和密码。
   - 调用服务层方法 userService.EmailLogin(u) 校验邮箱和密码的正确性：
     - 如果邮箱或密码错误，服务层返回错误，方法记录日志并返回错误信息 "Failed to login"。
4. **生成 JWT 令牌**：
   - 如果邮箱和密码验证成功，调用 `userApi.TokenNext(c, user)` 方法为登录用户生成一个 JWT 令牌。
   - 返回令牌和用户信息给客户端。

### QQ登录/QQLogin

```go
func (userApi *UserApi) QQLogin(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) QQLogin(accessTokenResponse other.AccessTokenResponse) (database.User, error)
```

#### 目的

`QQLogin` 方法的目的是实现基于 QQ 的第三方登录功能。通过 QQ 提供的 OAuth 授权机制，使用用户的授权码（`code`）获取访问令牌（`Access Token`）和用户唯一标识（`OpenID`），根据这些信息完成用户身份验证。如果登录成功，则生成 JWT 令牌并返回给客户端，用于后续的身份认证。

#### 流程

1. **获取授权码（`code`）**：
   - 从请求中获取 QQ 授权码 `code`（`c.Query("code")`）。
   - 判断 code 是否为空：
     - 如果为空，返回错误信息 "Code is required" 并终止流程。
2. **通过授权码获取访问令牌**：
   - 调用 `qqService.GetAccessTokenByCode(code)` 方法，将授权码发送给 QQ 的 OAuth 服务，获取 `Access Token` 和用户唯一标识 `OpenID`。
   - 判断返回的结果：
     - 如果出错（`err != nil`）或 `OpenID` 为空，记录错误日志并返回错误信息 "Invalid code"。
3. **根据 `Access Token` 和 `OpenID` 完成 QQ 登录**：
   - 调用服务层方法 `userService.QQLogin(accessTokenResponse)`，将 `AccessToken` 和 `OpenID` 传递给服务层，完成用户登录逻辑。
   - 服务层通常会根据 OpenID 检查用户是否已注册：
     - 如果用户已存在，直接返回用户信息。
     - 如果用户不存在，可能会自动为用户注册新账号（具体逻辑由服务层实现）。
   - 如果登录失败（`err != nil`），记录错误日志并返回错误信息 "Failed to login"。
4. **生成 JWT 令牌**：
   - 如果 QQ 登录成功，调用 `userApi.TokenNext(c, user)` 方法为用户生成一个 JWT 令牌。
   - 返回生成的令牌和用户信息给客户端。

### 生成并管理用户的访问令牌（JWT）和刷新令牌，同时通过 Redis 实现多点登录拦截/TokenNext

```go
func (userApi *UserApi) TokenNext(c *gin.Context, user database.User)
```

### 依赖

```go
func (jwtService *JwtService) SetRedisJWT(jwt string, uuid uuid.UUID) error
func (jwtService *JwtService) GetRedisJWT(uuid uuid.UUID) (string, error)
func (jwtService *JwtService) IsInBlacklist(jwt string) bool
func (jwtService *JwtService) JoinInBlacklist(jwtList database.JwtBlacklist) error
```

#### 目的

`TokenNext` 方法是登录逻辑的核心部分，主要用于生成并管理用户的访问令牌（JWT）和刷新令牌，同时通过 Redis 实现多点登录拦截，确保系统安全性和灵活性。
 此方法的目的包括：

1. **生成用户的访问令牌和刷新令牌**，实现用户认证与授权。
2. **处理多点登录拦截**，控制用户是否允许同时在多个设备登录。
3. **动态管理令牌状态**，支持令牌刷新和黑名单功能，防止不安全的令牌继续使用。

#### 流程

1. **检查用户状态**

- 检查传入的用户对象是否处于冻结状态（`user.Freeze`）。
- 如果用户被冻结，返回错误信息 "The user is frozen, contact the administrator" 并终止流程。

2. **生成访问令牌（Access Token）**

- 构造 `BaseClaims`（基础声明）对象，包含用户 ID（`UserID`）、UUID、角色 ID（`RoleID`）。
- 调用工具类方法 `CreateAccessToken` 生成一个有效期为 **2小时** 的访问令牌。
- 如果生成失败，记录错误日志并返回错误信息 "Failed to get accessToken"。

3. **生成刷新令牌（Refresh Token）**

- 基于相同的 `BaseClaims`，调用 `CreateRefreshToken` 生成一个有效期为 **7天** 的刷新令牌。
- 如果生成失败，记录错误日志并返回错误信息 "Failed to get refreshToken"。

4. **处理多点登录逻辑**

- 检查系统配置中是否启用了多点登录拦截（`global.Config.System.UseMultipoint`）。
- 未启用多点登录拦截：
  - 直接设置刷新令牌到客户端，并返回用户信息、访问令牌和过期时间。
- 启用多点登录拦截：
  - 检查 Redis 中是否已存在该用户的令牌：
    - 不存在旧令牌：
      - 将当前刷新令牌存入 Redis，绑定到用户的 UUID。
      - 设置刷新令牌到客户端，并返回登录成功的响应。
    - 存在旧令牌：
      - 将旧的令牌加入黑名单，确保旧令牌失效。
      - 将新的刷新令牌存入 Redis。
      - 设置刷新令牌到客户端，并返回登录成功的响应。

5. **错误处理**

- 如果在设置 Redis 或黑名单过程中出现错误，记录日志并返回相应的错误信息。

6. **返回用户登录结果**

- 返回用户信息、访问令牌（`AccessToken`）、访问令牌有效期（毫秒时间戳）以及其他相关信息。

### 通过邮箱验证码找回密码/ForgotPassword

```go
func (userApi *UserApi) ForgotPassword(c *gin.Context) 
```

#### 依赖

```go
func (userService *UserService) ForgotPassword(req request.ForgotPassword) error
```

#### 目的

`ForgotPassword` 方法实现了通过邮箱验证码找回密码的功能。用户提交邮箱、验证码和新密码后，系统会验证验证码的合法性和有效性，并将用户的密码更新为加密后的新密码。该方法旨在提供一种安全可靠的密码找回机制。

#### 流程

1. **接收并绑定请求参数**

- 定义一个 `request.ForgotPassword` 结构体，用于接收用户的请求参数（包括邮箱、验证码、新密码等）。
- 调用 `c.ShouldBindJSON(&req)` 将请求体中的 JSON 数据绑定到 `req` 中。
- 如果参数绑定失败（如缺少字段或格式错误），返回错误信息并终止流程。

2. **从会话中获取存储的邮箱和验证码**

- 获取会话对象 `session := sessions.Default(c)`。
- 通过会话存储的信息进行一系列验证：
  1. 邮箱一致性校验：
     - 获取会话中存储的邮箱 `savedEmail`。
     - 如果会话中的邮箱不存在或与用户提交的邮箱不一致，返回错误信息 "This email doesn't match the email to be verified" 并终止流程。
  2. 验证码一致性校验：
     - 获取会话中存储的验证码 `savedCode`。
     - 如果会话中的验证码不存在或与用户提交的验证码不一致，返回错误信息 "Invalid verification code" 并终止流程。
  3. 验证码有效性校验：
     - 获取会话中存储的验证码过期时间 `savedTime`。
     - 如果验证码已过期（当前时间大于过期时间），返回错误信息 "The verification code has expired, please resend it" 并终止流程。

3. **调用服务层更新用户密码**

- 调用服务层方法 userService.ForgotPassword(req)，将用户输入的新密码更新到数据库中：
  - 通常会对新密码进行加密存储（如使用 `bcrypt` 或类似算法）。
  - 更新操作成功返回 `nil`，失败返回错误信息（如用户不存在或数据库操作失败）。
- 如果更新失败，记录错误日志并返回错误信息 "Failed to retrieve the password"。

4. **返回成功响应**

- 如果上述操作全部成功，返回成功信息 "Successfully retrieved" 给客户端。

### 获取用户名片/UserCard

```go
func (userApi *UserApi) UserCard(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) UserCard(req request.UserCard) (response.UserCard, error) 
```

#### 目的

`UserCard` 方法用于获取用户的公开信息（用户名片），通过用户的唯一标识 `UUID` 查询其公开字段，如用户名、头像、所在地和个性签名等。 该方法的目的在于提供一个接口，允许其他用户或模块安全地访问用户的公开信息，而不会泄露敏感数据。

#### 流程

1. **接收并绑定查询参数**

- 从请求的查询参数中获取用户的 UUID：
  - 定义 `request.UserCard` 结构体，用于绑定 `UUID` 参数。
  - 调用 `c.ShouldBindQuery(&req)` 将查询参数绑定到 `req` 对象中。
- 如果参数绑定失败（如 `UUID` 参数缺失或格式错误），返回错误信息并终止流程。

2. **调用服务层查询用户公开信息**

- 调用服务层方法 `userService.UserCard(req)`，根据传入的 `UUID` 查询用户的公开信息。
- 服务层会根据 `UUID` 查询数据库，返回用户的公开信息（如用户名、头像、所在地、个性签名等）。
- 如果查询失败（如用户不存在或数据库操作错误），记录错误日志并返回错误信息 "Failed to get card"。

3. **返回用户的名片信息**

- 如果查询成功，将用户的名片信息（`userCard`）作为响应数据返回给客户端。

### 用户登出/Logout

```go
func (userApi *UserApi) Logout(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) Logout(c *gin.Context)
```

#### 目的

`Logout` 方法用于实现用户登出功能。它清除用户的登录状态，通过移除相关的认证信息（如 `Refresh Token`、JWT 令牌）来确保用户无法再继续使用当前会话，保证系统的安全性。该方法的主要目的是安全地注销用户，避免未授权的访问，清理登录相关的资源。

#### 流程

1. **调用服务层处理登出逻辑**

- 调用 userService.Logout(c) 方法，执行登出所需的一系列操作，包括：
  1. 清除 `Refresh Token` Cookie：
     - 从客户端的 Cookie 中移除 `Refresh Token`，使用户无法再使用该令牌刷新登录状态。
  2. 将 JWT 加入黑名单：
     - 将当前的 JWT 令牌加入黑名单，确保该令牌在有效期内也无法再被使用。
  3. 删除 Redis 中的令牌记录：
     - 如果启用了多点登录拦截（或登录状态管理），从 Redis 中删除与当前用户相关的令牌记录，彻底清除登录状态。

2. **返回登出成功的响应**

- 调用 `response.OkWithMessage("Successful logout", c)` 方法，返回成功的登出消息给客户端。

### 用户修改密码/UserResetPassword

```go
func (userApi *UserApi) UserResetPassword(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) UserResetPassword(req request.UserResetPassword) error
func (userService *UserService) Logout(c *gin.Context)
```

#### 目的

`UserResetPassword` 方法用于实现用户修改密码的功能。用户需要提供当前密码和新密码，系统会校验当前密码的正确性，然后将用户的新密码加密后存储到数据库中。
 该方法的主要目的是确保用户可以安全地更新自己的密码，同时通过旧密码校验防止恶意操作。

#### 流程

1. **从数据库中查询用户信息**

- 使用传入的 UserID 从数据库中查询用户记录：
  - 调用 `global.DB.Take(&user, req.UserID).Error` 查找用户。
  - 如果查询失败（如用户不存在），直接返回错误信息。

2. **校验当前密码的正确性**

- 调用 utils.BcryptCheck(req.Password, user.Password) 校验用户提供的当前密码是否与数据库中存储的密码一致：
  - 如果校验失败（密码不匹配），返回错误信息 `"original password does not match the current account"` 并终止流程。

3. **加密新密码**

- 如果密码校验成功，调用 `utils.BcryptHash(req.NewPassword)` 对用户提供的新密码进行加密处理。
- 将加密后的新密码赋值给用户对象的 `Password` 字段。

4. **更新数据库中的密码**

- 调用 `global.DB.Save(&user).Error` 方法，将更新后的用户密码存储到数据库中。
- 如果更新失败（如数据库操作错误），返回相应的错误信息。

5. **返回操作结果**

- 如果上述所有操作成功，返回 `nil` 表示密码重置成功。
- 如果出现任何错误（如用户不存在或数据库更新失败），返回错误信息。

### 获取当前登录用户的详细信息/UserInfo

```go
func (userApi *UserApi) UserInfo(c *gin.Context) 
```

#### 依赖

```go
func (userService *UserService) UserInfo(userID uint) (database.User, error) 
```

#### 目的

`UserInfo` 方法用于获取当前登录用户的详细信息。通过用户的唯一标识 `userID` 查询数据库，返回该用户的完整信息（如邮箱、角色权限、冻结状态、注册时间等）。
 该方法的主要目的是为用户提供其账户的详细信息，支持个人中心或账户管理功能。

#### 流程

1. **获取用户 ID**

- 调用工具方法 utils.GetUserID(c) 从当前会话上下文中获取用户的唯一标识 userID：
  - 该 `userID` 通常存储在 JWT 令牌中，解码后获取。
  - 如果未能获取到 `userID`，可能会引发后续操作失败。

2. **调用服务层查询用户信息**

- 调用服务层方法 userService.UserInfo(userID)，根据 userID查询数据库中对应用户的详细信息：
  - 服务层会根据 `userID` 从数据库中查找用户记录。
  - 如果查询失败（如用户不存在或数据库错误），返回错误信息。

3. **处理查询结果**

- 如果查询成功，将返回的用户信息对象（`user`）作为响应数据。
- 如果查询失败，记录错误日志并返回错误信息 "Failed to get user information"。

4. **返回用户信息**

- 调用 `response.OkWithData(user, c)` 方法，将查询到的用户信息以 JSON 格式返回给客户端。

### 更新当前登录用户的基本信息/UserChangeInfo

```go
func (userApi *UserApi) UserChangeInfo(c *gin.Context) 
```

#### 依赖

```go
func (userService *UserService) UserChangeInfo(req request.UserChangeInfo) error
```

#### 目的

`UserChangeInfo` 方法用于更新当前登录用户的基本信息（如头像、所在地、个性签名）。通过用户提交的请求数据，系统对指定字段进行验证并更新到数据库。
 该方法的主要目的是提供用户修改个人信息的能力，满足用户个性化设置需求，同时确保操作的安全性和有效性。

#### 流程

1. **接收并绑定请求参数**

- 定义一个 `request.UserChangeInfo` 结构体，用于接收客户端提交的修改信息（如头像、所在地、个性签名）。
- 调用 `c.ShouldBindJSON(&req)` 方法将请求体中的 JSON 数据绑定到 `req` 对象。
- 如果参数绑定失败（如缺少字段或格式错误），返回错误信息并终止流程。

2. **获取用户 ID**

- 调用工具方法 `utils.GetUserID(c)` 从当前会话上下文中提取用户的唯一标识 `UserID`。
- 将提取到的 `UserID` 赋值给 `req.UserID`，确保修改操作仅作用于当前登录用户。

3. **调用服务层更新用户信息**

- 调用服务层方法 userService.UserChangeInfo(req)，将修改后的用户信息更新到数据库：
  - 服务层根据 `UserID` 找到对应的用户记录。
  - 更新用户提交的字段（如头像、所在地、个性签名）。
- 如果更新失败（如用户不存在或数据库操作错误），记录错误日志并返回错误信息 "Failed to change user information"。

4. **返回成功响应**

- 如果所有操作成功，返回消息 "Successfully changed user information" 给客户端，表示用户信息更新完成。

### 获取用户地理位置的实时天气信息/UserWeather

```go
func (userApi *UserApi) UserWeather(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) UserWeather(ip string) (string, error)
```

#### 目的

`UserWeather` 方法通过用户的 IP 地址获取其地理位置，并查询高德天气 API 返回该地理位置的实时天气信息。 主要目的是提供基于用户位置的天气信息服务，提升用户体验，同时通过缓存优化性能。

#### 流程

1. **获取客户端 IP 地址**

- 调用 c.ClientIP() 方法获取用户的客户端 IP 地址：
  - 该 IP 地址用于定位用户的地理位置。
  - 如果获取 IP 失败，则可能导致后续操作失败。

2. **调用服务层获取天气信息**

- 调用 userService.UserWeather(ip) 方法，根据用户 IP 地址获取天气信息：
  1. 通过 IP 获取地理位置：
     - 使用 IP 定位服务（如高德 IP 定位 API）将用户 IP 转换为所在地理位置（如省份、城市）。
  2. 查询高德天气 API：
     - 根据定位到的地理位置（城市代码），调用高德天气 API 查询实时天气数据。
     - 返回天气信息，包括天气状况、温度、风向、风力等级、湿度等。
  3. 缓存结果：
     - 将获取到的天气信息缓存到 Redis 或其他缓存系统，设置缓存时间为 1 小时。
     - 如果缓存中存在相同 IP 的天气数据，则直接返回缓存结果而无需重复查询。

3. **错误处理**

- 如果服务层获取天气信息失败（如 IP 定位失败、天气 API 查询失败或缓存服务失败）：
  - 记录错误日志 `global.Log.Error`。
  - 返回错误信息 `"Failed to get user weather"` 给客户端。

4. **返回天气信息**

- 如果获取天气信息成功，将天气数据作为响应内容返回给客户端。

### 用于获取用户注册和登录的统计数据/UserChart

```go
func (userApi *UserApi) UserChart(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) UserChart(req request.UserChart) (response.UserChart, error) 
```

#### 目的

`UserChart` 方法用于获取用户注册和登录的统计数据（如每日注册数和每日登录数），并将其整理为图表所需的数据结构。此方法主要用于分析用户增长趋势，支持管理后台的可视化展示。

#### 流程

1. **接收并绑定查询参数**

- 定义一个 `request.UserChart` 结构体，用于接收查询参数（如统计的天数 `date`）。
- 调用 c.ShouldBindQuery(&req) 将请求的查询参数绑定到 req 对象中：
  - `date` 参数指定统计的天数（最大 30 天），如果未传则可能使用默认值。
- 如果参数绑定失败（如参数缺失或格式不正确），返回错误信息并终止流程。

2. **调用服务层获取统计数据**

- 调用服务层方法 userService.UserChart(req)，根据传入的统计天数获取用户的注册和登录数据：
  1. 查询注册和登录数据：
     - 服务层根据传入的天数（如最近 7 天或 30 天），查询数据库中每日的注册数和登录数。
     - 查询逻辑通常会按照日期分组聚合，生成每日的统计数据。
  2. 生成日期序列：
     - 服务层根据统计天数生成日期序列（如 `["2025-03-18", "2025-03-19", ...]`）。
     - 确保日期序列与统计结果对齐，填补没有数据的日期（如某天没有注册或登录记录时对应统计为 0）。
  3. 返回整理后的数据结构：
     - 返回一个包含日期序列、每日注册数和每日登录数的对象。

3. **处理服务层返回结果**

- 如果服务层返回的数据为空或发生错误：
  - 记录错误日志（`global.Log.Error`）。
  - 返回错误响应，提示用户 "Failed to get user chart"。
- 如果服务层返回数据成功，将结果整理后返回给客户端。

4. **返回图表数据**

- 调用 `response.OkWithData(data, c)` 方法，将整理后的统计数据以 JSON 格式返回给客户端。

### 管理员权限的接口，用于分页查询用户列表/UserList

```go
func (userApi *UserApi) UserList(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) UserList(info request.UserList) (interface{}, int64, error)
```

#### 目的

`UserList` 方法是管理员权限的接口，用于分页查询用户列表，并支持通过条件筛选（如用户 UUID、注册方式）获取用户信息。
 该方法的主要目的是为管理员提供用户管理功能，便于查看用户数据并执行相关操作（如分析用户注册来源）。

#### 流程

1. **接收并绑定查询参数**

- 定义一个 request.UserList 结构体，用于接收查询参数，包括分页信息和筛选条件：
  - 分页参数：
    - `page`：页码，默认值为 1。
    - `page_size`：每页数量，默认值为 10。
  - 筛选条件：
    - `uuid`：用户 UUID，用于精确筛选指定用户。
    - `register`：注册方式（1：邮箱，2：QQ）。
- 调用 `c.ShouldBindQuery(&pageInfo)` 方法将请求参数绑定到 `pageInfo` 对象中。
- 如果参数绑定失败（如缺少字段或格式错误），返回错误信息并终止流程。

2. **调用服务层查询用户列表**

- 调用服务层方法 `userService.UserList(pageInfo)`，根据分页参数和筛选条件查询用户数据：
  1. 分页查询：
     - 根据 `page` 和 `page_size` 参数，从数据库中取出指定页码的数据。
     - 计算偏移量 `offset = (page - 1) * page_size`，并限制查询数量为 `page_size`。
  2. 条件筛选：
     - 如果 `uuid` 参数不为空，则按用户 UUID 查询。
     - 如果 `register` 参数不为空，则按注册方式筛选用户（如邮箱注册或 QQ 注册）。
  3. 统计总数：
     - 在数据库中统计符合条件的用户总数 `total`，用于前端分页功能。
- 返回查询结果，包括用户列表 `list` 和总记录数 `total`。

3. **处理服务层返回结果**

- 如果服务层返回错误（如数据库查询失败），记录错误日志 `global.Log.Error`，并返回错误信息 `"Failed to get user list"`。
- 如果查询成功，将用户列表和总记录数封装为分页响应结构 `response.PageResult`。

4. **返回用户列表数据**

- 调用 `response.OkWithData` 方法，将分页响应数据以 JSON 格式返回给客户端。

### 管理员冻结普通用户账户/UserFreeze

```go
func (userApi *UserApi) UserFreeze(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) UserFreeze(req request.UserOperation) error
```

#### 目的

`UserFreeze` 方法用于管理员冻结普通用户账户。通过接收用户操作参数，校验管理员权限后更新用户的冻结状态，同时实现操作留痕和安全机制（如禁止冻结管理员账户）。该方法的主要目的是提供管理员对用户账户冻结的管理功能，确保系统的安全性和操作的合规性。

#### 流程

1. **接收并绑定请求参数**

- 定义一个 request.UserOperation

   结构体，用于接收管理员提交的操作请求：

  - 关键字段：
    - `UserID`：要冻结的用户 ID。
    - `Action`：操作类型（如冻结）。

- 使用 `c.ShouldBindJSON(&req)` 将请求体的 JSON 数据绑定到 `req` 对象中。

- 如果参数绑定失败（如字段缺失或格式错误），返回错误信息并终止流程。

2. **调用服务层执行冻结逻辑**

- 调用服务层方法 userService.UserFreeze(req)，根据传入的用户操作参数执行冻结逻辑：
  1. 校验管理员权限：
     - 确保当前操作人具有管理员权限。
     - 如果无权限，直接返回错误。
  2. 禁止冻结管理员账户：
     - 检查目标用户是否为管理员账户。
     - 如果目标账户为管理员，禁止冻结并返回错误信息。
  3. 更新用户冻结状态：
     - 在数据库中更新目标用户的冻结状态（如 `Freeze=true`）。
     - 冻结状态的变更通常涉及 `UPDATE` 操作，并可能记录操作时间。
  4. 记录操作日志：
     - 将操作记录存储到日志表中，包括操作人、操作时间、目标用户及操作内容，便于后续审计。

3. **处理服务层返回结果**

- 如果服务层返回错误（如权限不足、目标用户不存在、数据库操作失败）：
  - 记录错误日志 `global.Log.Error`。
  - 返回错误信息 `"Failed to freeze user"`。
- 如果冻结操作成功，返回成功消息 `"Successfully freeze user"`。

4. **返回操作结果**

- 调用 `response.OkWithMessage("Successfully freeze user", c)`，返回操作成功的消息给客户端。

### 管理员解冻普通用户账户/UserUnfreeze

```go
func (userApi *UserApi) UserUnfreeze(c *gin.Context)
```

#### 依赖

```go
func (userService *UserService) UserUnfreeze(req request.UserOperation) error
```

#### 目的

`UserUnfreeze` 方法用于管理员解冻被冻结的用户账户。通过接收用户操作参数，校验管理员权限后更新用户的冻结状态，同时记录操作日志并关联管理员 ID，确保操作可追溯性。
 该方法的主要目的是为管理员提供解冻用户账户的功能，以便恢复用户的使用权限，同时保持系统操作的安全性和可审计性。

#### 流程

1. **接收并绑定请求参数**

- 定义一个 request.UserOperation 结构体，用于接收管理员提交的操作请求：
  - 关键字段：
    - `UserID`：要解冻的用户 ID。
    - `Action`：操作类型（解冻）。
- 调用 `c.ShouldBindJSON(&req)` 方法将请求体中的 JSON 数据绑定到 `req` 对象。
- 如果参数绑定失败（如字段缺失或格式错误），返回错误信息并终止流程。

2. **调用服务层执行解冻逻辑**

- 调用服务层方法 userService.UserUnfreeze(req)，根据传入的用户操作参数执行解冻逻辑：
  1. 校验管理员权限：
     - 确保当前操作人具有管理员权限。
     - 如果无权限，直接返回错误信息。
  2. 检查用户账户状态：
     - 查询目标用户当前的冻结状态。
     - 如果目标用户未被冻结，则返回错误信息（如 "User is not frozen"）。
  3. 更新用户冻结状态：
     - 在数据库中更新目标用户的冻结状态为未冻结（如 `Freeze=false`）。
     - 通常会更新相关字段，如冻结状态和更新时间。
  4. 记录操作日志：
     - 将解冻操作记录到日志表中，包括操作管理员的 ID、时间、目标用户及操作内容。
     - 便于后续审计操作行为。

3. **处理服务层返回结果**

- 如果服务层返回错误（如权限不足、目标用户不存在或数据库操作失败）：
  - 记录错误日志 `global.Log.Error`。
  - 返回错误信息 `"Failed to unfreeze user"`。
- 如果解冻操作成功，返回成功消息 `"Successfully unfreeze user"`。

4. **返回操作结果**

- 调用 `response.OkWithMessage("Successfully unfreeze user", c)` 方法，返回操作成功的消息给客户端。

### 管理员获取用户的登录日志信息/UserLoginList

#### 依赖

```go
func (userService *UserService) UserLoginList(info request.UserLoginList) (interface{}, int64, error)
```

#### 目的

`UserLoginList` 方法用于获取用户的登录日志信息，支持分页查询和按用户 UUID 筛选。管理员可以通过此接口查看用户的登录时间、IP 地址、登录设备和登录地点等信息，用于系统管理和行为审计。该方法的主要目的是为管理员提供查询用户登录记录的能力，提升系统的安全性和管理效率。

#### 流程

1. **接收并绑定查询参数**

- 定义一个 request.UserLoginList 数据结构，用于接收查询参数，包括分页和筛选条件：
  - 分页参数：
    - `page`：页码，默认值为 1（可选）。
    - `page_size`：每页数量，默认值为 10（可选）。
  - 筛选条件：
    - `uuid`：用户 UUID，用于查询指定用户的登录记录（可选）。
- 调用 `c.ShouldBindQuery(&pageInfo)` 方法，将请求中的查询参数绑定到 `pageInfo` 对象。
- 如果参数绑定失败（如字段缺失或格式错误），返回错误信息并终止流程。

2. **调用服务层查询登录日志**

- 调用服务层方法 userService.UserLoginList(pageInfo)，根据分页参数和筛选条件查询用户的登录日志：
  1. 分页查询：
     - 根据 `page` 和 `page_size` 参数，从数据库中查询指定页码的数据。
     - 计算偏移量 `offset = (page - 1) * page_size`，并限制查询数量为 `page_size`。
  2. 条件筛选：
     - 如果 `uuid` 参数不为空，则按用户 UUID 筛选登录记录。
  3. 返回结果：
     - 返回登录日志列表 `list` 和符合条件的总记录数 `total`。

3. **处理服务层返回结果**

- 如果服务层返回错误（如数据库查询失败），记录错误日志 `global.Log.Error`，并返回错误响应 `"Failed to get user login list"`。
- 如果查询成功，将登录日志列表和总记录数封装到分页响应结构 `response.PageResult` 中。

4. **返回分页结果**

- 调用 `response.OkWithData(response.PageResult, c)` 方法，将登录日志数据以 JSON 格式返回给客户端。