# DuckDB Portal

## 简介
DuckDB Portal 是一个基于DuckDB的Web界面工具，旨在简化数据库管理和数据分析流程。

## 功能特性
- 支持多种数据源导入导出
- 提供SQL查询编辑器
- 可视化数据展示

## 安装步骤
1. 克隆仓库
   ```bash
   git clone https://github.com/xqbumu/duckdb-portal.git
   cd duckdb-portal
   ```
2. 安装依赖
   ```bash
   go mod download
   ```
3. 启动服务
   ```bash
   go run main.go
   ```

## 使用方法
1. 打开浏览器访问 `http://localhost:8080`
2. 登录系统
3. 使用SQL编辑器执行查询
4. 查看查询结果

## 贡献
欢迎提交问题和Pull Request！

// 添加 GitHub Actions 使用 goreleaser
## GitHub Actions
我们使用 GitHub Actions 来自动化版本发布流程。具体配置文件位于 `.github/workflows/release.yml`。