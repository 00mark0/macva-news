import React, { useState } from "react";

const Counter = () => {
	const [count, setCount] = useState(0);

	if (count < 0) {
		setCount(0);
	}

	return (
		<div>
			<h2 className="text-red-300">React Counter</h2>
			<p>Count: {count}</p>
			<button className="text-green-500" onClick={() => setCount(count + 1)}>Increment</button>
			<button className="text-red-500" onClick={() => setCount(count - 1)}>Decrement</button>
		</div>
	);
};

export default Counter;

