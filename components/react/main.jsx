import React from 'react';
import { createRoot } from 'react-dom/client';
import TrendingContent from './trendingContent';
import DailyAnalytics from './daily-analytics';

// Registry mapping component names to React components
const components = {
	TrendingContent,
	DailyAnalytics
};

// Keep track of mounted components to prevent duplicate mounts
const mountedElements = new WeakSet();

// Function to mount React components in a given root element
export default function mountReactComponents(root = document) {
	root.querySelectorAll('[data-react-component]').forEach((el) => {
		// Skip if already mounted
		if (mountedElements.has(el)) return;

		const componentName = el.getAttribute('data-react-component');
		const Component = components[componentName];

		if (Component) {
			// Mark as mounted before rendering to prevent race conditions
			mountedElements.add(el);

			// Priority mounting for visible components
			if (isElementInViewport(el) || componentName === 'TrendingContent') {
				createRoot(el).render(<Component />);
			} else {
				// Defer mounting of off-screen components
				requestIdleCallback(() => {
					createRoot(el).render(<Component />);
				}, { timeout: 2000 });
			}
		} else {
			console.error(`No component registered for ${componentName}`);
		}
	});
}

// Helper to check if element is in viewport
function isElementInViewport(el) {
	const rect = el.getBoundingClientRect();
	return (
		rect.top >= 0 &&
		rect.left >= 0 &&
		rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
		rect.right <= (window.innerWidth || document.documentElement.clientWidth)
	);
}

// Run initial mounting ASAP for critical components
if (document.readyState === 'loading') {
	document.addEventListener('DOMContentLoaded', () => {
		mountReactComponents();

		// Listen for HTMX swaps and mount React inside new content
		document.body.addEventListener('htmx:afterSwap', (event) => {
			mountReactComponents(event.target);
		});
	});
} else {
	// Document already loaded, mount immediately
	mountReactComponents();

	// Listen for HTMX swaps
	document.body.addEventListener('htmx:afterSwap', (event) => {
		mountReactComponents(event.target);
	});
}

// Polyfill for requestIdleCallback
if (!window.requestIdleCallback) {
	window.requestIdleCallback = (callback, options) => {
		const timeout = options?.timeout || 50;
		return setTimeout(() => {
			callback({
				didTimeout: false,
				timeRemaining: () => 0
			});
		}, timeout);
	};

	window.cancelIdleCallback = (id) => {
		clearTimeout(id);
	};
}

