import React from 'react';
import { createRoot } from 'react-dom/client';
import Counter from './counter';
import Widget from './widget';

// A registry mapping component names to React components
const components = {
	Counter,
	Widget,
};

// Find all elements with a data-react-component attribute
document.querySelectorAll('[data-react-component]').forEach((el) => {
	const componentName = el.getAttribute('data-react-component');
	const Component = components[componentName];
	if (Component) {
		// Optionally, you can pass props by reading data from the element
		createRoot(el).render(<Component />);
	} else {
		console.error(`No component registered for ${componentName}`);
	}
});

