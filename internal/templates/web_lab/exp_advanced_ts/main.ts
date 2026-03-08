interface Box<T> {
  value: T;
}

const payload: Box<number> = { value: 42 };
console.log(payload.value);
