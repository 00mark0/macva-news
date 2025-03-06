import { useState, useEffect } from 'react'
import { format, startOfDay, startOfWeek, subDays } from 'date-fns'
import fetch from './axios';

export default function TrendingContent() {
	const [trendingContent, setTrendingContent] = useState([]);
	const [loading, setLoading] = useState(false);
	const [limit, setLimit] = useState(12);
	const [dateFilter, setDateFilter] = useState('day');
	const [hasMore, setHasMore] = useState(true);
	const [error, setError] = useState(null);

	/*function getCookie(name) {
		const value = `; ${document.cookie}`;
		const parts = value.split(`; ${name}=`);
		if (parts.length === 2) return parts.pop().split(';').shift();
		return null;
	}*/

	//const token = getCookie('access_token')

	const getPublishedAtDate = () => {
		const now = new Date();
		let date;

		switch (dateFilter) {
			case 'day':
				date = startOfDay(now); // Start of current day
				break;
			case 'week':
				date = startOfWeek(now, { weekStartsOn: 1 }); // Start of current week
				break;
			case 'month':
				date = subDays(now, 30); // Start of current month
				break;
			default:
				date = startOfDay(now);
		}

		const formattedDate = format(date, "yyyy-MM-dd");

		return formattedDate;
	};

	const fetchTrendingContent = async () => {
		setLoading(true);
		setError(null);
		try {
			const publishedAt = getPublishedAtDate();
			const res = await fetch.get(`/api/admin/trending?published_at=${publishedAt}&limit=${limit}`);

			// Handle null response by setting to empty array
			if (res.data === null) {
				setTrendingContent([]);
				setHasMore(false);
			} else {
				setTrendingContent(res.data);
				setHasMore(res.data.length === limit);
			}

		} catch (error) {
			console.error('Error fetching trending content:', error);
			setError('Failed to fetch trending content');
			// Ensure we have an empty array in case of error
			setTrendingContent([]);
		} finally {
			setLoading(false);
		}
	}

	const handleLoadMore = () => {
		setLimit(prevLimit => prevLimit + 12);
	}

	const handleDateFilterChange = (filter) => {
		setDateFilter(filter);
		setLimit(12);
		// Reset content state when changing filters to prevent flash of old content
		setTrendingContent([]);
	}

	useEffect(() => {
		fetchTrendingContent();
	}, [limit, dateFilter]);

	const getFilterLabel = () => {
		switch (dateFilter) {
			case 'day': return 'danas';
			case 'week': return 'ove nedelje';
			case 'month': return 'ovog meseca';
			default: return 'danas';
		}
	};

	// Safety check to ensure trendingContent is always an array
	const safeContent = Array.isArray(trendingContent) ? trendingContent : [];

	return (
		<div className="w-full dark:bg-black sm:p-8 p-4">
			<h1 className="text-black dark:text-white text-2xl font-bold mb-6">Trending sadržaj</h1>

			{/* Date filter buttons */}
			<div className="flex space-x-4 mb-6">
				<button
					onClick={() => handleDateFilterChange('day')}
					className={`px-4 py-2 rounded-md ${dateFilter === 'day'
						? 'bg-blue-600 text-white'
						: 'bg-gray-200 dark:bg-gray-700 text-black dark:text-white'
						}`}
				>
					Danas
				</button>
				<button
					onClick={() => handleDateFilterChange('week')}
					className={`px-4 py-2 rounded-md ${dateFilter === 'week'
						? 'bg-blue-600 text-white'
						: 'bg-gray-200 dark:bg-gray-700 text-black dark:text-white'
						}`}
				>
					Ove nedelje
				</button>
				<button
					onClick={() => handleDateFilterChange('month')}
					className={`px-4 py-2 rounded-md ${dateFilter === 'month'
						? 'bg-blue-600 text-white'
						: 'bg-gray-200 dark:bg-gray-700 text-black dark:text-white'
						}`}
				>
					Poslednjih 30 dana
				</button>
			</div>

			{/* Loading state */}
			{loading && safeContent.length === 0 ? (
				<div className="flex justify-center py-8">
					<div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			) : (
				<>
					{/* Error message */}
					{error && (
						<div className="bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-200 p-4 rounded-md mb-6">
							{error}. Molimo pokušajte ponovo kasnije.
						</div>
					)}

					{/* No content message */}
					{!loading && !error && safeContent.length === 0 && (
						<div className="bg-gray-100 dark:bg-gray-800 p-6 rounded-lg text-center mb-6">
							<svg className="w-16 h-16 mx-auto text-gray-400 dark:text-gray-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
								<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"></path>
							</svg>
							<h3 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-2">
								Nema trending sadržaja {getFilterLabel()}
							</h3>
							<p className="text-gray-600 dark:text-gray-400 mb-4">
								Pokušajte drugi vremenski period ili se vratite kasnije.
							</p>
							<div className="flex justify-center space-x-3">
								{dateFilter !== 'week' && (
									<button
										onClick={() => handleDateFilterChange('week')}
										className="cursor-pointer px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
									>
										Probaj ove nedelje
									</button>
								)}
								{dateFilter !== 'month' && (
									<button
										onClick={() => handleDateFilterChange('month')}
										className="cursor-pointer px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
									>
										Probaj poslednjih 30 dana
									</button>
								)}
							</div>
						</div>
					)}

					{/* Content list */}
					{safeContent.length > 0 && (
						<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
							{safeContent.map((item) => (
								<div
									key={item.content_id}
									className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden transition-all hover:shadow-lg flex flex-col hover:scale-105"
								>
									<div className="p-4 flex flex-col flex-grow">
										<h2 className="h-12 text-lg font-semibold text-black dark:text-white mb-2 truncate">{item.title}</h2>
										<p className="text-gray-600 dark:text-gray-300 text-sm mb-4 line-clamp-3 flex-grow">
											{item.content_description}
										</p>

										{/* Content stats */}
										<div className="flex flex-wrap overflow-auto items-center justify-between text-sm text-gray-500 dark:text-gray-400 mt-auto">
											<div className="flex space-x-4">
												<span className="flex items-center">
													<svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
														<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
														<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
													</svg>
													{item.view_count}
												</span>

												<span className="flex items-center">
													<svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
														<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905C11 5.37 10.5 7 9 8.5 7.5 10 5.5 10 4 10h-.5"></path>
													</svg>
													{item.like_count}
												</span>

												<span className="flex items-center">
													<svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
														<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 14h-4.764a2 2 0 00-1.789-2.894l3.5-7A2 2 0 008.737 3h4.017c.163 0 .326.02.485.06l4.764 3.88m-7 10V19a2 2 0 002 2h.095c.5 0 .905-.405.905-.905C13 18.63 13.5 17 15 15.5 16.5 14 18.5 14 20 14h.5"></path>
													</svg>
													{item.dislike_count}
												</span>

												<span className="flex items-center">
													<svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
														<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"></path>
													</svg>
													{item.comment_count}
												</span>
											</div>

											{/* Displaying date inline */}
											<span className="text-sm">
												{new Date(item.published_at).toLocaleDateString('en-GB')}
											</span>
										</div>
									</div>
								</div>
							))}
						</div>
					)}
				</>
			)}

			{/* Loading indicator for more content */}
			{loading && safeContent.length > 0 && (
				<div className="flex justify-center py-6">
					<div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			)}

			{/* Load more button */}
			{!loading && hasMore && safeContent.length > 0 && (
				<div className="flex justify-center mt-8">
					<button
						onClick={handleLoadMore}
						className="cursor-pointer px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
					>
						Učitaj više
					</button>
				</div>
			)}
		</div>
	);
}
