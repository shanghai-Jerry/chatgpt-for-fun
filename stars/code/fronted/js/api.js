// API接口调用封装

// 基础URL配置
const BASE_URL = 'http://localhost:8080';

// API端点
const API_ENDPOINTS = {
    GOALS: '/goals',
    GOAL_BY_ID: (id) => `/goals/${id}`,
    GOALS_BY_CATEGORY: (category) => `/goals/category/${category}`,
    TOTAL_STARS: '/stars',
    DAILY_RATING: (id) => `/goals/${id}/daily-rating`,
    DAILY_RATINGS: (id) => `/goals/${id}/daily-ratings`,
    // 评论相关端点
    COMMENTS: (id) => `/goals/${id}/comments`
};

// HTTP请求工具函数
const http = {
    // GET请求
    get: async (url) => {
        try {
            const response = await fetch(BASE_URL + url);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error('GET请求失败:', error);
            throw error;
        }
    },
    
    // POST请求
    post: async (url, data) => {
        try {
            const response = await fetch(BASE_URL + url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error('POST请求失败:', error);
            throw error;
        }
    },
    
    // PUT请求
    put: async (url, data) => {
        try {
            const response = await fetch(BASE_URL + url, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error('PUT请求失败:', error);
            throw error;
        }
    },
    
    // DELETE请求
    delete: async (url) => {
        try {
            const response = await fetch(BASE_URL + url, {
                method: 'DELETE'
            });
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error('DELETE请求失败:', error);
            throw error;
        }
    }
};

// 目标相关API
const goalAPI = {
    // 获取所有目标
    getAllGoals: () => http.get(API_ENDPOINTS.GOALS),
    
    // 根据ID获取目标
    getGoalById: (id) => http.get(API_ENDPOINTS.GOAL_BY_ID(id)),
    
    // 根据类别获取目标
    getGoalsByCategory: (category) => http.get(API_ENDPOINTS.GOALS_BY_CATEGORY(category)),
    
    // 创建新目标
    createGoal: (goalData) => http.post(API_ENDPOINTS.GOALS, goalData),
    
    // 更新目标
    updateGoal: (id, goalData) => http.put(API_ENDPOINTS.GOAL_BY_ID(id), goalData),
    
    // 删除目标
    deleteGoal: (id) => http.delete(API_ENDPOINTS.GOAL_BY_ID(id)),
    
    // 获取总星数
    getTotalStars: () => http.get(API_ENDPOINTS.TOTAL_STARS),
    
    // 为指定目标添加每日评分
    addDailyRating: (goalId, ratingData) => http.post(API_ENDPOINTS.DAILY_RATING(goalId), ratingData),
    
    // 获取指定目标的所有每日评分记录
    getDailyRatings: (goalId) => http.get(API_ENDPOINTS.DAILY_RATINGS(goalId)),
    
    // 为指定目标创建评论
    createComment: (goalId, commentData) => http.post(API_ENDPOINTS.COMMENTS(goalId), commentData),
    
    // 获取指定目标的所有评论
    getComments: (goalId) => http.get(API_ENDPOINTS.COMMENTS(goalId))
};