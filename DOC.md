项目结构
---

添加一个 API endpoint 需要以下步骤：
1. 在 model/syzoj.api.proto 里添加对应的 Page，然后用 go generate 生成代码。
2. 在前端的 components 里编写页面并在 components/any.js 添加对应的 Page。
3. 在后端的 server/handlers 里编写后端并在 server/handlers/handlers.go 添加对应的 Route.
