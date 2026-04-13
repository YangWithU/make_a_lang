# 从 0 到 1：可执行版 TODO（Go 语言解释器/可扩展编译器）

> 目标：初学者按本文逐条实现，最终得到**可运行用户程序**的语言实现。
>
> 路线：先解释器（MVP 必做），后字节码编译器+VM（进阶选做）。

---

## 0. 先统一术语（避免后续混乱）

- 本项目第一阶段交付物是**解释器**（Lexer + Parser + AST + Evaluator）。
- 你原始需求写“编译器”，这里做工程化修正：
  - MVP：树遍历解释器（可运行脚本、支持闭包/作用域/控制流）。
  - Advanced：Compiler + VM（字节码执行）。

---

## 1. 最终交付标准（Definition of Done）

必须同时满足：

- [ ] `go test ./...` 全绿。
- [ ] `go run ./cmd/lang` 可进入 REPL。
- [ ] `go run ./cmd/lang run examples/smoke.lang` 正确执行并输出预期。
- [ ] 支持：整数、布尔、null、字符串、变量声明、重赋值、优先级表达式、if/else、while、函数、return、闭包、print/input。
- [ ] 对类型错误/未定义变量/参数不匹配给出可读错误并中止当前程序执行。

---

## 2. 目录结构与里程碑

### 2.1 建立项目骨架（阶段 0）

- [ ] 执行：`go mod init make_a_lang`
- [ ] 创建目录：
  - [ ] `cmd/lang`
  - [ ] `internal/token`
  - [ ] `internal/lexer`
  - [ ] `internal/ast`
  - [ ] `internal/parser`
  - [ ] `internal/object`
  - [ ] `internal/environment`
  - [ ] `internal/evaluator`
  - [ ] `internal/builtin`
  - [ ] `internal/repl`
  - [ ] `examples`

**阶段验收**
- [ ] `cmd/lang/main.go` 先打印 `lang bootstrap ok`，保证工程能启动。

---

## 3. 阶段 1：Token 设计 + Lexer（可直接照写）

这一阶段回答你问的核心问题：

- 如何设计字符读取状态？
- 先写什么 struct？
- Token 如何抽象？
- Tokenizer/Lexer 怎样落地？

### 3.1 先实现 Token 抽象（必须先做）

#### 文件
- [ ] `internal/token/token.go`

#### 你要写的类型（先把框架搭出来）
- [ ] `type TokenType string`
- [ ] `type Token struct { Type TokenType; Literal string; Line int; Column int }`

#### 你要定义的 token 常量（MVP）
- [ ] 特殊：`ILLEGAL`, `EOF`
- [ ] 标识符与字面量：`IDENT`, `INT`, `STRING`
- [ ] 运算：`ASSIGN`, `PLUS`, `MINUS`, `ASTERISK`, `SLASH`, `BANG`, `EQ`, `NOT_EQ`, `LT`, `GT`
- [ ] 分隔符：`COMMA`, `SEMICOLON`, `LPAREN`, `RPAREN`, `LBRACE`, `RBRACE`
- [ ] 关键字：`FUNCTION`, `LET`, `TRUE`, `FALSE`, `IF`, `ELSE`, `RETURN`, `WHILE`, `NULL`

#### 你要写的函数
- [ ] `func LookupIdent(ident string) TokenType`
  - 作用：判断标识符是否关键字，不是则返回 `IDENT`。

#### 自测
- [ ] 写表驱动单测：`LookupIdent("fn") == FUNCTION`，`LookupIdent("x") == IDENT`。

---

### 3.2 实现 Lexer 字符读取状态（你问的重点）

#### 文件
- [ ] `internal/lexer/lexer.go`

#### 先写 Lexer struct（建议直接用 byte 版本，先不处理 Unicode 复杂度）

- [ ] 
```go
type Lexer struct {
    input        string
    position     int  // 当前字符下标（指向 ch）
    readPosition int  // 下一个要读的字符下标
    ch           byte // 当前字符

    line         int  // 当前行号，从 1 开始
    column       int  // 当前列号，从 1 开始
}
```

#### 为什么这样设计
- [ ] `position/readPosition/ch` 是经典“单字符前瞻”模型，足够处理 `==`、`!=`。
- [ ] `line/column` 让错误定位可读。

#### 你要实现的函数顺序（照顺序写）

1) [ ] `func New(input string) *Lexer`
- 初始化 `line=1, column=0`。
- 调一次 `readChar()` 读入首字符。

2) [ ] `func (l *Lexer) readChar()`
- 如果 `readPosition >= len(input)`，设置 `ch=0` 表示 EOF。
- 否则 `ch = input[readPosition]`。
- `position = readPosition`，`readPosition++`。
- 处理位置信息：
  - 若读到 `\n`：`line++`，`column=0`。
  - 否则 `column++`。

3) [ ] `func (l *Lexer) peekChar() byte`
- 返回下一个字符但不前进；越界返回 `0`。

4) [ ] `func (l *Lexer) skipWhitespace()`
- 连续跳过 ` `、`\t`、`\n`、`\r`。

5) [ ] `func (l *Lexer) skipComment()`
- MVP 只支持 `//` 单行注释。
- 当当前是 `/` 且 `peekChar() == '/'`，读到换行或 EOF。

6) [ ] `func (l *Lexer) readIdentifier() string`
- 从当前下标开始读 `[a-zA-Z_]` 开头，后续允许 `[a-zA-Z0-9_]`。

7) [ ] `func (l *Lexer) readNumber() string`
- 连续读取数字字符。

8) [ ] `func (l *Lexer) readString() (string, error)`
- 假设字符串由双引号包裹。
- 当前字符是 `"` 时进入，读取直到下一个 `"` 或 EOF。
- MVP 可先不支持转义；若 EOF 未闭合，返回错误。

9) [ ] `func (l *Lexer) NextToken() token.Token`
- 这是 Tokenizer 的主入口：每次调用返回一个 token。
- 流程固定：
  - `skipWhitespace()`
  - 如果是注释，先 `skipComment()` 再继续跳过空白
  - `switch l.ch` 分支处理单字符与双字符 token
  - 标识符/数字/字符串分别走对应读取函数
  - 生成 token 时填入 `Line/Column`
  - 在合适时机调用 `readChar()` 前进

#### Tokenizer vs Lexer（术语说明）
- [ ] 在这个项目里可以视为同一层：`Lexer.NextToken()` 就是 tokenizer 行为。

#### 这一阶段的最低测试集合（必须写）

- [ ] `internal/lexer/lexer_test.go` 增加：
  - [ ] 基础 token 序列测试（覆盖 let/fn/if/while/return/操作符/分隔符）。
  - [ ] 双字符操作符测试（`==`, `!=`）。
  - [ ] 字符串测试（含空格）。
  - [ ] 注释跳过测试。
  - [ ] 行列号测试（至少验证 2~3 个 token 位置）。

**阶段验收**
- [ ] 给下面输入，token 序列必须稳定：

```txt
let x = 10;
// c1
if (x == 10) { x = x + 1; }
print("ok");
```

---

## 4. 阶段 2：AST + Pratt Parser（可执行拆解）

### 4.1 AST 节点定义顺序

#### 文件
- [ ] `internal/ast/ast.go`

#### 先定义接口
- [ ] `Node`：`TokenLiteral() string`、`String() string`
- [ ] `Statement`：嵌入 `Node` + `statementNode()`
- [ ] `Expression`：嵌入 `Node` + `expressionNode()`

#### 按顺序实现节点（不要跳）
- [ ] `Program`
- [ ] `Identifier`
- [ ] `LetStatement`
- [ ] `AssignStatement`（重赋值必须独立语句节点）
- [ ] `ReturnStatement`
- [ ] `ExpressionStatement`
- [ ] `IntegerLiteral`, `StringLiteral`, `BooleanLiteral`, `NullLiteral`
- [ ] `PrefixExpression`, `InfixExpression`
- [ ] `BlockStatement`
- [ ] `IfExpression`
- [ ] `WhileStatement`
- [ ] `FunctionLiteral`, `CallExpression`

**AST 自测**
- [ ] `Program.String()` 输出可读（用于 parser 测试对比）。

---

### 4.2 Parser 结构与 Pratt 入口

#### 文件
- [ ] `internal/parser/parser.go`

#### 先定义 Parser struct
- [ ]
```go
type Parser struct {
    l         *lexer.Lexer
    curToken  token.Token
    peekToken token.Token
    errors    []string

    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns  map[token.TokenType]infixParseFn
}
```

#### 你要先实现的辅助函数
- [ ] `nextToken()`
- [ ] `curTokenIs()` / `peekTokenIs()`
- [ ] `expectPeek()`（失败时写入错误）
- [ ] `peekError()`

#### 定义优先级
- [ ] `LOWEST, EQUALS, LESSGREATER, SUM, PRODUCT, PREFIX, CALL`
- [ ] `precedences map[token.TokenType]int`

#### 实现语句解析入口
- [ ] `ParseProgram()`
- [ ] `parseStatement()` 分派：
  - [ ] `parseLetStatement()`
  - [ ] `parseReturnStatement()`
  - [ ] `parseWhileStatement()`
  - [ ] `parseAssignStatement()`（判定逻辑：`cur=IDENT && peek=ASSIGN`）
  - [ ] 默认走 `parseExpressionStatement()`

#### 实现表达式解析（Pratt）
- [ ] `parseExpression(precedence int)`
- [ ] 注册 prefix：ident/int/string/bool/null/`!`/`-`/group/if/function
- [ ] 注册 infix：`+ - * / == != < >` 与 call

#### 错误恢复（必须做）
- [ ] 单条语句解析失败后，跳到 `;` 或 `}` 再继续。
- [ ] `Errors()` 对外暴露 parser 错误。

#### 阶段测试（必须）
- [ ] let/return/assign/while 语句解析测试。
- [ ] 运算符优先级测试（字符串化 AST 对比）。
- [ ] 函数定义与调用解析测试。
- [ ] parser 错误收集测试。

---

## 5. 阶段 3：对象系统与环境链

#### 文件
- [ ] `internal/object/object.go`
- [ ] `internal/environment/environment.go`

### 5.1 Object 设计

- [ ] `type ObjectType string`
- [ ] `type Object interface { Type() ObjectType; Inspect() string }`
- [ ] 实现对象：
  - [ ] `Integer`
  - [ ] `Boolean`
  - [ ] `String`
  - [ ] `Null`
  - [ ] `ReturnValue`
  - [ ] `Error`
  - [ ] `Function{Parameters, Body, Env}`
  - [ ] `Builtin{Fn func(args ...Object) Object}`

### 5.2 Environment 设计（闭包关键）

- [ ] struct：
```go
type Environment struct {
    store map[string]object.Object
    outer *Environment
}
```
- [ ] 构造函数：`New()` 与 `NewEnclosed(outer *Environment)`
- [ ] `Get(name)`：当前找不到则递归 outer
- [ ] `Set(name, val)`：绑定当前作用域
- [ ] `Assign(name, val)`：沿作用域链找已存在变量并更新；找不到返回错误

#### 阶段测试
- [ ] 覆盖嵌套环境查找、遮蔽（shadowing）、外层赋值更新。

---

## 6. 阶段 4：Evaluator（执行引擎）

#### 文件
- [ ] `internal/evaluator/evaluator.go`

#### 核心函数
- [ ] `func Eval(node ast.Node, env *environment.Environment) object.Object`

#### 实现顺序（务必按顺序）

1) [ ] 字面量求值（int/bool/string/null）
2) [ ] 前缀表达式（`!`, `-`）
3) [ ] 中缀表达式（算术、比较、字符串拼接）
4) [ ] 程序与块求值（遇 `ReturnValue`/`Error` 立即短路返回）
5) [ ] let 与 assign
6) [ ] if/else
7) [ ] while
8) [ ] 函数对象、调用、参数绑定
9) [ ] 闭包（函数定义时捕获 env）
10) [ ] 内置函数表（print/input）

#### 关键规则（必须写进代码）

- [ ] truthy 规则：仅 `false` 与 `null` 为假，其余为真。
- [ ] 类型错误统一构造 `Error` 对象，不 panic。
- [ ] `return` 仅在函数边界解包。
- [ ] 函数参数数量不匹配时报错。

#### 阶段测试（必须）
- [ ] 算术与优先级
- [ ] 比较与逻辑非
- [ ] if/else 与 while
- [ ] 函数调用与 return
- [ ] 闭包 `makeAdder`
- [ ] 错误传播（`1 + "x"`、未定义变量）

---

## 7. 阶段 5：Builtins、CLI、REPL

#### 文件
- [ ] `internal/builtin/builtin.go`
- [ ] `internal/repl/repl.go`
- [ ] `cmd/lang/main.go`

### 7.1 Builtins
- [ ] `print(args ...Object)`：打印并返回 `null`
- [ ] `input()`：读取 stdin 一行，返回 `String`

### 7.2 REPL
- [ ] 循环读一行 -> lex -> parse -> eval -> 打印结果
- [ ] REPL 使用同一个全局 env（跨行保留变量）
- [ ] parser 错误与 runtime 错误分开输出

### 7.3 CLI
- [ ] 无参数：REPL
- [ ] `run <file>`：读取脚本并执行

#### 阶段验收
- [ ] `go run ./cmd/lang` 可交互运行。
- [ ] `go run ./cmd/lang run examples/smoke.lang` 可执行脚本。

---

## 8. 阶段 6：测试与样例程序（交付前必须）

### 测试组织
- [ ] `internal/lexer/lexer_test.go`
- [ ] `internal/parser/parser_test.go`
- [ ] `internal/evaluator/evaluator_test.go`

### 样例脚本
- [ ] `examples/smoke.lang`
- [ ] `examples/closure.lang`
- [ ] `examples/loop.lang`
- [ ] `examples/io.lang`

### 冒烟脚本（手动）
- [ ] 连续执行：
  - [ ] `go test ./...`
  - [ ] `go run ./cmd/lang run examples/smoke.lang`

---

## 9. 阶段 7：文档与新人可复现性

- [ ] 新建 `README.md`，至少包含：
  - [ ] 语法支持列表
  - [ ] 运行方式（REPL / run file）
  - [ ] 常见错误说明
  - [ ] 开发路线图（解释器 -> VM）
- [ ] 新建 `docs/grammar.md`（可选但强烈建议）
  - [ ] 写出 MVP 语法（EBNF 或伪 BNF）

---

## 10. 进阶：编译器 + VM（可选）

- [ ] 新建 `internal/compiler`：AST -> Bytecode
- [ ] 新建 `internal/vm`：执行字节码
- [ ] 指令最小集：
  - [ ] 常量加载
  - [ ] 算术与比较
  - [ ] 跳转
  - [ ] 调用与返回
- [ ] 跑同一批 evaluator 测试样例，比较输出一致性

---

## 11. 每天执行模版（保证不拖延）

每天按下面四步走：

1. [ ] 先实现 1~2 个小函数（不跨模块）
2. [ ] 马上补对应单测
3. [ ] 跑 `go test ./...`
4. [ ] 在 TODO 打勾并写一句“今天完成了什么”

---

## 12. 最小验收程序（最终必须通过）

```javascript
let x = 1 + 2 * 3;
print(x); // 7

let makeAdder = fn(a) {
  return fn(b) { return a + b; };
};

let add2 = makeAdder(2);
print(add2(5)); // 7

let i = 0;
while (i < 3) {
  print(i);
  i = i + 1;
}

print(input());
```

预期：输出 `7, 7, 0, 1, 2`，并回显输入。若出现 `1 + "x"`，必须报类型错误并停止。

---

## 13. 常见实现坑（提前规避）

- [ ] 坑 1：在 `NextToken()` 里遗漏 `readChar()`，导致死循环。
- [ ] 坑 2：字符串读取没处理 EOF，导致越界或无限循环。
- [ ] 坑 3：赋值 `x = 1` 误当作声明处理，破坏作用域语义。
- [ ] 坑 4：`return` 未短路传播，函数体后续代码仍继续执行。
- [ ] 坑 5：闭包调用时没用定义时环境，导致外层变量丢失。

---

## 14. 你现在就可以开始的第一组任务（今天）

- [ ] 完成 `internal/token/token.go`
- [ ] 完成 `internal/lexer/lexer.go` 的 `Lexer` struct + `readChar` + `peekChar`
- [ ] 完成 `skipWhitespace` + `readIdentifier` + `readNumber`
- [ ] 在 `NextToken()` 支持：`let x = 10;` 所需 token
- [ ] 写第一条 lexer 单测并跑通

只要完成这 5 项，你就进入“可持续推进”状态了。
