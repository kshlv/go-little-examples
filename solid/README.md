# SOLID

I wonder why haven't I thought of it earlier?

SOLID is actually an abbreviate and stands for those five principles:
* Single Responsibility Principle
* Open / Closed Principle
* Liskov Substitution Principle
* Interface Segregation Principle
* Dependency Invertion Principle

Let's take a look at each one.

## Single Responsibility Principle

_SRP_ sounds something like this:

> Each class must response for just one actor.

For example, by making this we violate _SRP_:
```go
type Student struct {
	firstTermHomeWork, firstTermTest, firstTermPaper    float64
	secondTermHomeWork, secondTermTest, secondTermPaper float64
}

// FirstTermGrade ...
func (s *Student) FirstTermGrade() float64 {
	return (s.firstTermHomeWork + s.firstTermTest + s.firstTermPaper) / 3.0
}

// SecondTermGrade ...
func (s *Student) SecondTermGrade() float64 {
	return (s.secondTermHomeWork + s.secondTermTest + s.secondTermPaper) / 3.0
}
```

The reason for that is that the class (_whatever_) is responsible for calculating both of the terms. 
So, if we what to apply changes to this class, there might be _two_ reasons for that, and that goes against _SRP_.

We can change it by doing something like this:
```go
type Grade struct {
	Name                  string
	homework, test, paper float64
}

func (g *Grade) Grade() float64 {
	return (g.homework + g.test + g.paper) / 3.0
}

type Student struct {
	First, Second *Grade
}

```

Now, if we want to change the grade calculation, these's the only point of responsibility.

## Open / Closed Principle

The _OCP_ sounds like:

> One software entity must be open for extension but closed for modification.

So, once a class implements a certain scope of requirements, the further implementation sould not need to change in proder to fulfill requirements.

```go
type MyLogger struct {
    format string
}

func New() *MyLogger {
    return &MyLogger{
        format: "%s: %s\n",
    }
}

func (logger *MyLogger) Log(msg ...interface{}) {
	fmt.Fprintf(os.Stdout, logger.format, time.Now().Unix(), msg)
}
```

But once we change the format string to, say, `"[log] %s: %s\n"`, we break the backward compability as well as the _OCP_.

This is what we should do instead:
```go
func New(format string) *MyLogger {
    return &MyLogger{
        format: format,
    }
}
```

# Barbara Liskov Principle

```go
type Person interface {
    Greet() string
}

type Student interface {
    Person
    YearsOld(age int) string
}

```

