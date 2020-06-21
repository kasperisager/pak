// https://github.com/estree/estree
package ast

type (
	Program struct {
		SourceType string
		Body       []ProgramBody
	}

	ProgramBody interface {
		VisitProgramBody(ProgramBodyVisitor)
	}

	ProgramBodyVisitor struct {
		Statement         func(Statement)
		ModuleDeclaration func(ModuleDeclaration)
	}

	Expression interface {
		VisitExpression(ExpressionVisitor)
	}

	ExpressionVisitor struct {
		Identifier               func(*Identifier)
		Literal                  func(Literal)
		ThisExpression           func(*ThisExpression)
		ArrayExpression          func(*ArrayExpression)
		ObjectExpression         func(*ObjectExpression)
		FunctionExpression       func(*FunctionExpression)
		UnaryExpression          func(*UnaryExpression)
		UpdateExpression         func(*UpdateExpression)
		BinaryExpression         func(*BinaryExpression)
		AssignmentExpression     func(*AssignmentExpression)
		LogicalExpression        func(*LogicalExpression)
		ConditionalExpression    func(*ConditionalExpression)
		CallExpression           func(*CallExpression)
		NewExpression            func(*NewExpression)
		SequenceExpression       func(*SequenceExpression)
		ArrowFunctionExpression  func(*ArrowFunctionExpression)
		YieldExpression          func(*YieldExpression)
		AwaitExpression          func(*AwaitExpression)
		TemplateLiteral          func(*TemplateLiteral)
		TaggedTemplateExpression func(*TaggedTemplateExpression)
		ClassExpression          func(*ClassExpression)
	}

	Statement interface {
		ProgramBody
		VisitStatement(StatementVisitor)
	}

	StatementVisitor struct {
		Declaration         func(Declaration)
		ExpressionStatement func(*ExpressionStatement)
		BlockStatement      func(*BlockStatement)
		EmptyStatement      func(*EmptyStatement)
		DebuggerStatement   func(*DebuggerStatement)
		WithStatement       func(*WithStatement)
		ReturnStatement     func(*ReturnStatement)
		LabeledStatement    func(*LabeledStatement)
		BreakStatement      func(*BreakStatement)
		ContinueStatement   func(*ContinueStatement)
		IfStatement         func(*IfStatement)
		SwitchStatement     func(*SwitchStatement)
		ThrowStatement      func(*ThrowStatement)
		TryStatement        func(*TryStatement)
		WhileStatement      func(*WhileStatement)
		DoWhileStatement    func(*DoWhileStatement)
		ForStatement        func(*ForStatement)
		ForInStatement      func(*ForInStatement)
		ForOfStatement      func(*ForOfStatement)
	}

	Declaration interface {
		Statement
		VisitDeclaration(DeclarationVisitor)
	}

	DeclarationVisitor struct {
		FunctionDeclaration func(*FunctionDeclaration)
		VariableDeclaration func(*VariableDeclaration)
		ClassDeclaration    func(*ClassDeclaration)
	}

	Pattern interface {
		VisitPattern(PatternVisitor)
	}

	PatternVisitor struct {
		Identifier       func(*Identifier)
		MemberExpression func(*MemberExpression)
	}

	Identifier struct {
		Name string
	}

	Literal interface {
		Expression
		VisitLiteral(LiteralVisitor)
	}

	LiteralVisitor struct {
		StringLiteral  func(*StringLiteral)
		BooleanLiteral func(*BooleanLiteral)
		NullLiteral    func(*NullLiteral)
		NumberLiteral  func(*NumberLiteral)
		RegExpLiteral  func(*RegExpLiteral)
	}

	StringLiteral struct {
		Value string
	}

	BooleanLiteral struct {
		Value bool
	}

	NullLiteral struct{}

	NumberLiteral struct {
		Value float64
	}

	RegExpLiteral struct {
		Regex struct {
			Pattern string
			Flags   string
		}
	}

	ExpressionStatement struct {
		Expression Expression
	}

	BlockStatement struct {
		Body []Statement
	}

	EmptyStatement struct{}

	DebuggerStatement struct{}

	WithStatement struct {
		Object Expression
		Body   Statement
	}

	ReturnStatement struct {
		Argument Expression
	}

	LabeledStatement struct {
		Label *Identifier
		Body  Statement
	}

	BreakStatement struct {
		Label *Identifier
	}

	ContinueStatement struct {
		Label *Identifier
	}

	IfStatement struct {
		Test       Expression
		Consequent Statement
		Alternate  Statement
	}

	SwitchStatement struct {
		Discriminant Expression
		Cases        []*SwitchCase
	}

	SwitchCase struct {
		Test       Expression
		Consequent []Statement
	}

	ThrowStatement struct {
		Argument Expression
	}

	TryStatement struct {
		Block     *BlockStatement
		Handler   *CatchClause
		Finalizer *BlockStatement
	}

	CatchClause struct {
		Param Pattern
		Body  *BlockStatement
	}

	WhileStatement struct {
		Test Expression
		Body Statement
	}

	DoWhileStatement struct {
		Test Expression
		Body Statement
	}

	ForStatement struct {
		Init   ForStatementInit
		Test   Expression
		Update Expression
		Body   Statement
	}

	ForStatementInit interface {
		VisitForStatementInit(ForStatementInitVisitor)
	}

	ForStatementInitVisitor struct {
		VariableDeclaration func(*VariableDeclaration)
		Expression          func(Expression)
	}

	ForInStatement struct {
		Left  ForInStatementLeft
		Right Expression
		Body  Statement
	}

	ForInStatementLeft interface {
		VisitForInStatementLeft(ForInStatementLeftVisitor)
	}

	ForInStatementLeftVisitor struct {
		VariableDeclaration func(*VariableDeclaration)
		Pattern             func(Pattern)
	}

	ForOfStatement struct {
		Left  ForOfStatementLeft
		Right Expression
		Body  Statement
		Await bool
	}

	ForOfStatementLeft interface {
		VisitForOfStatementLeft(ForOfStatementLeftVisitor)
	}

	ForOfStatementLeftVisitor struct {
		VariableDeclaration func(*VariableDeclaration)
		Pattern             func(Pattern)
	}

	FunctionDeclaration struct {
		ID        *Identifier
		Params    []Pattern
		Body      *BlockStatement
		Generator bool
		Async     bool
	}

	VariableDeclaration struct {
		Kind         string
		Declarations []*VariableDeclarator
	}

	VariableDeclarator struct {
		ID   Pattern
		Init Expression
	}

	ThisExpression struct{}

	ArrayExpression struct {
		Elements []ArrayExpressionElement
	}

	ArrayExpressionElement interface {
		VisitArrayExpressionElement(ArrayExpressionElementVisitor)
	}

	ArrayExpressionElementVisitor struct {
		Expression    func(Expression)
		SpreadElement func(*SpreadElement)
	}

	ObjectExpression struct {
		Properties []ObjectExpressionProperty
	}

	ObjectExpressionProperty interface {
		VisitObjectExpressionProperty(ObjectExpressionPropertyVisitor)
	}

	ObjectExpressionPropertyVisitor struct {
		Property      func(*Property)
		SpreadElement func(*SpreadElement)
	}

	Property struct {
		Key       Expression
		Value     Expression
		Kind      string
		Method    bool
		Shorthand bool
		Computed  bool
	}

	FunctionExpression struct {
		ID        *Identifier
		Params    []Pattern
		Body      *BlockStatement
		Generator bool
		Async     bool
	}

	UnaryExpression struct {
		Operator string
		Prefix   bool
		Argument Expression
	}

	UpdateExpression struct {
		Operator string
		Prefix   bool
		Argument Expression
	}

	BinaryExpression struct {
		Operator string
		Left     Expression
		Right    Expression
	}

	AssignmentExpression struct {
		Operator string
		Left     Pattern
		Right    Expression
	}

	LogicalExpression struct {
		Operator string
		Left     Expression
		Right    Expression
	}

	MemberExpression struct {
		Object   MemberExpressionObject
		Property Expression
		Computed bool
	}

	MemberExpressionObject interface {
		VisitMemberExpressionObject(MemberExpressionObjectVisitor)
	}

	MemberExpressionObjectVisitor struct {
		Expression func(Expression)
		Super      func(*Super)
	}

	ConditionalExpression struct {
		Test       Expression
		Alternate  Expression
		Consequent Expression
	}

	CallExpression struct {
		Callee    CallExpressionCallee
		Arguments []CallExpressionArgument
	}

	CallExpressionCallee interface {
		VisitCallExpressionCallee(CallExpressionCalleeVisitor)
	}

	CallExpressionCalleeVisitor struct {
		Expression func(Expression)
		Super      func(*Super)
	}

	CallExpressionArgument interface {
		VisitCallExpressionArgument(CallExpressionArgumentVisitor)
	}

	CallExpressionArgumentVisitor struct {
		Expression    func(Expression)
		SpreadElement func(*SpreadElement)
	}

	NewExpression struct {
		Callee    Expression
		Arguments []NewExpressionArgument
	}

	NewExpressionArgument interface {
		VisitNewExpressionArgument(NewExpressionArgumentVisitor)
	}

	NewExpressionArgumentVisitor struct {
		Expression    func(Expression)
		SpreadElement func(*SpreadElement)
	}

	SequenceExpression struct {
		Expression []Expression
	}

	ArrowFunctionExpression struct {
		Body       ArrowFunctionExpressionBody
		Expression bool
	}

	ArrowFunctionExpressionBody interface {
		VisitArrowFunctionExpressionBody(ArrowFunctionExpressionBodyVisitor)
	}

	ArrowFunctionExpressionBodyVisitor struct {
		BlockStatement func(*BlockStatement)
		Expression     func(Expression)
	}

	YieldExpression struct {
		Argument Expression
		Delegate bool
	}

	AwaitExpression struct {
		Argument Expression
	}

	TemplateLiteral struct {
		Quasis      []*TemplateElement
		Expressions []Expression
	}

	TaggedTemplateExpression struct {
		Tag   Expression
		Quasi TemplateLiteral
	}

	TemplateElement struct {
		Tail  bool
		Value struct {
			Cooked *string
			Raw    string
		}
	}

	ObjectPattern struct {
		Properties []ObjectPatternProperty
	}

	ObjectPatternProperty interface {
		VisitObjectPatternProperty(ObjectPatternPropertyVisitor)
	}

	ObjectPatternPropertyVisitor struct {
		AssignmentProperty func(*AssignmentProperty)
		RestElement        func(*RestElement)
	}

	AssignmentProperty struct {
		Key       Expression
		Value     Pattern
		Shorthand bool
		Computed  bool
	}

	ArrayPattern struct {
		Elements []Pattern
	}

	RestElement struct {
		Argument Pattern
	}

	AssignmentPattern struct {
		Left  Pattern
		Right Expression
	}

	Super struct{}

	SpreadElement struct {
		Argument Expression
	}

	Class struct {
		ID         *Identifier
		SuperClass Expression
		Body       ClassBody
	}

	ClassBody struct {
		Body []*MethodDefinition
	}

	MethodDefinition struct {
		Key      Expression
		Value    *FunctionExpression
		Kind     string
		Computed bool
		Static   bool
	}

	ClassDeclaration struct {
		ID         *Identifier
		SuperClass Expression
		Body       *ClassBody
	}

	ClassExpression struct {
		ID         *Identifier
		SuperClass Expression
		Body       *ClassBody
	}

	MetaProperty struct {
		Meta     *Identifier
		Property *Identifier
	}

	ModuleDeclaration interface {
		VisitModuleDeclaration(ModuleDeclarationVisitor)
	}

	ModuleDeclarationVisitor struct {
		ImportDeclaration        func(*ImportDeclaration)
		ExportNamedDeclaration   func(*ExportNamedDeclaration)
		ExportDefaultDeclaration func(*ExportDefaultDeclaration)
		ExportAllDeclaration     func(*ExportAllDeclaration)
	}

	ImportDeclaration struct {
		Specifiers []ImportDeclarationSpecifier
		Source     *StringLiteral
	}

	ImportDeclarationSpecifier interface {
		VisitImportDeclarationSpecifier(ImportDeclarationSpecifierVisitor)
	}

	ImportDeclarationSpecifierVisitor struct {
		ImportSpecifier          func(*ImportSpecifier)
		ImportDefaultSpecifier   func(*ImportDefaultSpecifier)
		ImportNamespaceSpecifier func(*ImportNamespaceSpecifier)
	}

	ImportSpecifier struct {
		Local    *Identifier
		Imported *Identifier
	}

	ImportDefaultSpecifier struct {
		Local *Identifier
	}

	ImportNamespaceSpecifier struct {
		Local *Identifier
	}

	ExportNamedDeclaration struct {
		Declaration Declaration
		Specifiers  []*ExportSpecifier
		Source      *StringLiteral
	}

	ExportSpecifier struct {
		Local    *Identifier
		Exported *Identifier
	}

	ExportDefaultDeclaration struct {
		Declaration ExportDefaultDeclarationDeclaration
	}

	ExportDefaultDeclarationDeclaration interface {
		VisitExportDefaultDeclarationDeclaration(ExportDefaultDeclarationDeclarationVisitor)
	}

	ExportDefaultDeclarationDeclarationVisitor struct {
		AnonymousDefaultExportedFunctionDeclaration func(*AnonymousDefaultExportedFunctionDeclaration)
		FunctionDeclaration                         func(*FunctionDeclaration)
		AnonymousDefaultExportedClassDeclaration    func(*AnonymousDefaultExportedClassDeclaration)
		ClassDeclaration                            func(*ClassDeclaration)
		Expression                                  func(Expression)
	}

	AnonymousDefaultExportedFunctionDeclaration struct {
		Params    []Pattern
		Body      *BlockStatement
		Generator bool
	}

	AnonymousDefaultExportedClassDeclaration struct {
		SuperClass Expression
		Body       *ClassBody
	}

	ExportAllDeclaration struct {
		Source Literal
	}
)

func (p *ExpressionStatement) VisitProgramBody(v ProgramBodyVisitor) { v.Statement(p) }
func (p *BlockStatement) VisitProgramBody(v ProgramBodyVisitor)      { v.Statement(p) }
func (p *ImportDeclaration) VisitProgramBody(v ProgramBodyVisitor)   { v.ModuleDeclaration(p) }

func (s *ExpressionStatement) VisitStatement(v StatementVisitor) { v.ExpressionStatement(s) }
func (s *BlockStatement) VisitStatement(v StatementVisitor)      { v.BlockStatement(s) }

func (e *Identifier) VisitExpression(v ExpressionVisitor)            { v.Identifier(e) }
func (e *StringLiteral) VisitExpression(v ExpressionVisitor)         { v.Literal(e) }
func (e *BooleanLiteral) VisitExpression(v ExpressionVisitor)        { v.Literal(e) }
func (e *NullLiteral) VisitExpression(v ExpressionVisitor)           { v.Literal(e) }
func (e *NumberLiteral) VisitExpression(v ExpressionVisitor)         { v.Literal(e) }
func (e *RegExpLiteral) VisitExpression(v ExpressionVisitor)         { v.Literal(e) }
func (e *SequenceExpression) VisitExpression(v ExpressionVisitor)    { v.SequenceExpression(e) }
func (e *AssignmentExpression) VisitExpression(v ExpressionVisitor)  { v.AssignmentExpression(e) }
func (e *ConditionalExpression) VisitExpression(v ExpressionVisitor) { v.ConditionalExpression(e) }
func (e *LogicalExpression) VisitExpression(v ExpressionVisitor)     { v.LogicalExpression(e) }
func (e *BinaryExpression) VisitExpression(v ExpressionVisitor)      { v.BinaryExpression(e) }
func (e *UnaryExpression) VisitExpression(v ExpressionVisitor)       { v.UnaryExpression(e) }
func (e *UpdateExpression) VisitExpression(v ExpressionVisitor)      { v.UpdateExpression(e) }

func (l *StringLiteral) VisitLiteral(v LiteralVisitor)  { v.StringLiteral(l) }
func (l *BooleanLiteral) VisitLiteral(v LiteralVisitor) { v.BooleanLiteral(l) }
func (l *NullLiteral) VisitLiteral(v LiteralVisitor)    { v.NullLiteral(l) }
func (l *NumberLiteral) VisitLiteral(v LiteralVisitor)  { v.NumberLiteral(l) }
func (l *RegExpLiteral) VisitLiteral(v LiteralVisitor)  { v.RegExpLiteral(l) }

func (m *ImportDeclaration) VisitModuleDeclaration(v ModuleDeclarationVisitor) { v.ImportDeclaration(m) }

func (i *ImportSpecifier) VisitImportDeclarationSpecifier(v ImportDeclarationSpecifierVisitor) {
	v.ImportSpecifier(i)
}
func (i *ImportDefaultSpecifier) VisitImportDeclarationSpecifier(v ImportDeclarationSpecifierVisitor) {
	v.ImportDefaultSpecifier(i)
}
func (i *ImportNamespaceSpecifier) VisitImportDeclarationSpecifier(v ImportDeclarationSpecifierVisitor) {
	v.ImportNamespaceSpecifier(i)
}

func (p *Identifier) VisitPattern(v PatternVisitor)       { v.Identifier(p) }
func (p *MemberExpression) VisitPattern(v PatternVisitor) { v.MemberExpression(p) }
