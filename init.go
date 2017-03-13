package main

const initData = `(def inc (fn
  (n)
  (+ n 1)))

(def dec (fn
  (n)
  (- n 1)))

(def each (fn (f l)
    (if l
       (do (f (head l))
           (recur f (tail l))))))

(def filter (fn (f l)
  (_filter f l '())))

(def _filter (fn (f l r)
    (if l
       (if (f (head l))
         (recur f (tail l) (cons (head l) r))
         (recur f (tail l) r))
       r)))

(def map (fn (f l)
  (_map3 f l '())))

(def _map3 (fn (f l r)
  (if l
    (recur f (tail l) (cons (f (head l)) r))
    r)))

(def reduce (fn (f s l)
    (if l
      (recur f (f s (head l)) (tail l))
      s)))`
