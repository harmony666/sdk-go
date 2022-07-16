<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"API Breaking" for breaking exported APIs used by developers building on SDK.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased]

### Features

### Improvements

### API Breaking Changes

### Bug Fixes

### Deprecated

## v2.2.0 - 2021-12-17

### Features

* (gas) 新增启用/停用gas计费开关API
* (gas) 新增 attach `Limit` API
* (通用) 新增修改地址类型API
* (通用) 提供至信链地址生成相关API
* (Grpc client) grpc客户端发送消息时，可设置允许单条message大小的最大值(MB)

### Improvements

* (订阅) 支持订阅断线自动重连机制
* (订阅) 支持合约事件按照区块高度订阅历史事件
