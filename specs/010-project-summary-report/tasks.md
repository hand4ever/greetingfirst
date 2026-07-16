# 任务(Tasks): Greeting 项目阶段总结报告

**输入(Input)**: Design documents from `/specs/010-project-summary-report/`

**前置条件(Prerequisites)**: plan.md（必需）、spec.md（已完成）

**测试(Tests)**: 纯文档报告，不包含测试任务。

**组织(Organization)**: 任务按报告生成和交付阶段组织。

## 格式(Format): `[ID] [P?] [Story] 描述(Description)`

- **[P]**: 可并行执行 (不同文件，无依赖)
- **[Story]**: 不属于用户故事（纯文档报告，无 US 划分）
- 描述中包含具体文件路径

## 路径约定(Path Conventions)

- 报告源文件: `specs/010-project-summary-report/spec.md`
- PDF 输出: `specs/010-project-summary-report/Greeting项目阶段总结报告.pdf`

---

## 阶段(Phase) 1: 数据核验(Data Verification)

**目的(Purpose)**: 确保报告中所有量化数据准确可复验

- [ ] T001 重新运行 `git log --oneline | wc -l` 验证 commit 数量
- [ ] T002 [P] 重新运行 `find . -name "*.go" ! -name "*_test.go" ! -path "./.git/*" -exec cat {} + | wc -l` 验证 Go 生产代码行数
- [ ] T003 [P] 重新运行 `find . -name "*_test.go" ! -path "./.git/*" -exec cat {} + | wc -l` 验证测试代码行数
- [ ] T004 [P] 重新运行 `find ./specs -name "*.md" -exec cat {} + | wc -l` 验证 Spec 文档行数
- [ ] T005 [P] 对照各 spec 目录下 tasks.md 重新统计 Spec 001–009 任务完成度

**检查点(Checkpoint)**: 所有数据已核验，与报告一致

---

## 阶段(Phase) 2: 报告整理(Finalize Report)

**目的(Purpose)**: 完善报告内容，准备 PDF 导出

- [ ] T006 检查 spec.md 中表格格式、排版、链接完整性
- [ ] T007 根据核验结果更新 spec.md 中任何需要修正的数据

**检查点(Checkpoint)**: 报告内容最终确认

---

## 阶段(Phase) 3: PDF 导出(PDF Export)

**目的(Purpose)**: 将 Markdown 报告转换为格式精美的 PDF 文件

- [ ] T008 使用 PDF skill 将 `specs/010-project-summary-report/spec.md` 转换为高质量的 PDF 文件，输出路径 `specs/010-project-summary-report/Greeting项目阶段总结报告.pdf`，要求包含：
  - 封面页（项目名称、报告标题、日期）
  - 目录自动生成
  - 表格带样式（交替行颜色、边框）
  - 章节标题层次分明
  - 页眉/页脚（含页码）

**检查点(Checkpoint)**: PDF 文件生成完毕，可直接用于汇报

---

## 依赖与执行顺序(Dependencies & Execution Order)

### 阶段依赖(Phase Dependencies)

- **数据核验（阶段 1）**: 无依赖 - 可立即开始
- **报告整理（阶段 2）**: 依赖阶段 1 完成（需要核验结果）
- **PDF 导出（阶段 3）**: 依赖阶段 2 完成（需要最终版 spec.md）

### 并行机会(Parallel Opportunities)

- T002、T003、T004、T005 可并行执行（不同统计命令，无依赖）
- T001 也可与其他核验任务并行

---

## 实施策略(Implementation Strategy)

### 一次性完成

1. 完成阶段 1：并行运行所有数据核验命令
2. 完成阶段 2：修正并最终确认报告内容
3. 完成阶段 3：生成高质量 PDF

### 预估工作量

| 阶段 | 任务数 | 预估时间 |
|------|:--:|------|
| 阶段 1: 数据核验 | 5 | ~2 分钟 |
| 阶段 2: 报告整理 | 2 | ~3 分钟 |
| 阶段 3: PDF 导出 | 1 | ~5 分钟 |
| **合计** | **8** | **~10 分钟** |

---

## 备注(Notes)

- 本 spec 为文档型报告，无代码变更，无用户故事
- 最终交付物：`specs/010-project-summary-report/Greeting项目阶段总结报告.pdf`
- PDF 需富含格式美化（表格、封面、目录、页眉页脚）
