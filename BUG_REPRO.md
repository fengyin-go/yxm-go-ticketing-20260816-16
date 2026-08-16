# BUG_REPRO

## Bug 是什么
创建重复分类或重复邮箱客服时，底层冲突错误在向上传递时被 `%v` 重新包装，丢失了 `ErrConflict` 哨兵错误，导致 HTTP 层无法正确映射为 409。

## 如何触发
1. 创建分类 `A`，再次用相同名称创建分类。
2. 创建邮箱为 `dup@example.com` 的客服，再次用相同邮箱创建客服。

## 错误信息
接口返回 `500`，而不是 `409 Conflict`；`errors.Is(err, store.ErrConflict)` 为 `false`。
