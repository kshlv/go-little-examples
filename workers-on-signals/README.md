# Workers on os.Signals

Should `Ctrl+C` or `Ctrl+Z` always shut you running program down? No, not necessarily.

Should `SIGKILL` always kill you program? [Yeah, it should](http://www.cs.kent.edu/~ruttan/sysprog/lectures/signals.html#catching). There are two non-catchable signals: `SIGKILL` and `SIGSTOP`. You cannot catch and handle them in any normal means.

Can you do stuff to you program with `kill` other that actually killing it? It looks like yeah, lots of stuff.

Anyways, the goal of this exercise is to try and make something small and funny which would work on `os.Signal`s. Workers factory looks like a nice idea, so here it is. Also, it's always nice to do something conscious with channels so that not to forget how they work.
