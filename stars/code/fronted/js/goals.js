// 目标管理功能实现

// 渲染星星评分控件
// @param {number} goalId - 目标ID
// @param {number} currentStars - 当前星数
function renderStarRating(goalId, currentStars) {
    let starsHTML = '';
    for (let i = 1; i <= 5; i++) {
        // 根据当前星数决定是否添加filled类
        const starClass = i <= currentStars ? 'star filled' : 'star';
        starsHTML += `<span class="${starClass}" data-goal-id="${goalId}" data-star="${i}">⭐</span>`;
    }
    return starsHTML;
}

// 绑定星星点击事件
function bindStarClickEvents() {
    // 为所有星星绑定点击事件
    document.querySelectorAll('.star').forEach(star => {
        star.addEventListener('click', function() {
            const goalId = parseInt(this.getAttribute('data-goal-id'));
            const starValue = parseInt(this.getAttribute('data-star'));
            updateGoalStars(goalId, starValue);
        });
    });
}

// 加载总星数
async function loadTotalStars() {
    try {
        const data = await goalAPI.getTotalStars();
        document.getElementById('total-stars-value').textContent = data.total_stars;
    } catch (error) {
        console.error('加载总星数失败:', error);
        // 如果加载失败，显示为0
        document.getElementById('total-stars-value').textContent = '0';
    }
}

// 加载所有目标
async function loadAllGoals() {
    try {
        const goals = await goalAPI.getAllGoals();
        renderGoalList(goals);
        // 同时更新总星数
        loadTotalStars();
        // 绑定星星点击事件
        bindStarClickEvents();
    } catch (error) {
        console.error('加载目标列表失败:', error);
        document.getElementById('goal-list-container').innerHTML = 
            '<div class="no-goals">加载目标列表失败，请稍后重试</div>';
    }
}

// 根据类别加载目标
async function loadGoalsByCategory(category) {
    try {
        const goals = await goalAPI.getGoalsByCategory(category);
        renderGoalList(goals);
        // 同时更新总星数
        loadTotalStars();
        // 绑定星星点击事件
        bindStarClickEvents();
    } catch (error) {
        console.error('加载目标列表失败:', error);
        document.getElementById('goal-list-container').innerHTML = 
            '<div class="no-goals">加载目标列表失败，请稍后重试</div>';
    }
}

// 更新目标星数
async function updateGoalStars(goalId, stars) {
    try {
        // 创建每日评分记录
        const today = new Date();
        today.setHours(0, 0, 0, 0); // 设置时间为当天的开始
        
        const ratingData = {
            goal_id: goalId,
            rating: stars,
            date: today.toISOString(),
        };
        
        // 添加每日评分记录
        await goalAPI.addDailyRating(goalId, ratingData);
        
        // 重新加载目标列表
        await loadAllGoals();
        
        // 直接调用loadTotalStars更新总星数
        await loadTotalStars();
        
        // 更新首页的总星数
        if (window.location.pathname.includes('goal-list.html')) {
            // 如果在目标列表页面，通知首页更新总星数
            window.postMessage({type: 'UPDATE_TOTAL_STARS'}, '*');
        }
    } catch (error) {
        console.error('更新目标星数失败:', error);
        alert('更新目标星数失败，请稍后重试');
    }
}

// 渲染目标详情
function renderGoalDetail(goal) {
    const container = document.getElementById('goal-detail-container');
    
    if (!goal) {
        container.innerHTML = '<div class="no-goals">未找到目标</div>';
        return;
    }
    
    // 修改renderGoalDetail函数中的detailHTML，移除评星交互相关代码
    const detailHTML = `
        <h2>${goal.title}</h2>
        <span class="category">${goal.category}</span>
        <div class="stars">⭐ ${goal.stars} 星</div>
        <div class="description">${goal.description || '暂无描述'}</div>
        <div class="meta">
            <p>创建时间: ${new Date(goal.created_at).toLocaleString()}</p>
            <p>更新时间: ${new Date(goal.updated_at).toLocaleString()}</p>
        </div>
        <div class="daily-ratings" id="daily-ratings-${goal.id}">
            <h3>每日评分记录</h3>
            <div class="loading">加载中...</div>
        </div>
    `;
    
    container.innerHTML = detailHTML;
    
    // 移除绑定星星点击事件的调用
    // bindStarClickEvents();
    
    // 加载并显示每日评分记录
    loadDailyRatings(goal.id);
}

// 渲染目标列表
function renderGoalList(goals) {
    const container = document.getElementById('goal-list-container');
    
    if (!goals || goals.length === 0) {
        container.innerHTML = '<div class="no-goals">暂无目标</div>';
        return;
    }
    
    // 生成目标列表HTML
    const goalsHTML = goals.map(goal => `
        <div class="goal-item" data-id="${goal.id}">
            <h3>${goal.title}</h3>
            <span class="category">${goal.category}</span>
            <div class="stars">⭐ ${goal.stars} 星</div>
            <div class="star-rating-list">
                ${renderStarRating(goal.id, goal.stars)}
            </div>
            <p class="description">${goal.description || '暂无描述'}</p>
            <div class="goal-actions">
                <a href="goal-detail.html?id=${goal.id}" class="btn btn-small">查看详情</a>
            </div>
        </div>
    `).join('');
    
    container.innerHTML = goalsHTML;
}

// 加载并显示目标详情
async function loadGoalDetail(goalId) {
    try {
        // 调用API获取目标详情
        const goal = await goalAPI.getGoalById(goalId);
        
        // 渲染目标详情
        renderGoalDetail(goal);
    } catch (error) {
        console.error('加载目标详情失败:', error);
        document.getElementById('goal-detail-container').innerHTML = 
            '<div class="no-goals">加载目标详情失败，请稍后重试</div>';
    }
}

// 加载并显示每日评分记录
async function loadDailyRatings(goalId) {
    try {
        const ratings = await goalAPI.getDailyRatings(goalId);
        const container = document.getElementById(`daily-ratings-${goalId}`);
        
        if (!ratings || ratings.length === 0) {
            container.innerHTML = '<p>暂无每日评分记录</p>';
            return;
        }
        
        // 按日期排序
        ratings.sort((a, b) => new Date(b.date) - new Date(a.date));
        
        // 生成评分记录HTML
        const ratingsHTML = ratings.map(rating => `
            <div class="rating-item">
                <span class="rating-date">${new Date(rating.date).toLocaleDateString()}</span>
                <span class="rating-stars">${'⭐'.repeat(rating.rating)}</span>
            </div>
        `).join('');
        
        container.innerHTML = `
            <h3>每日评分记录</h3>
            <div class="ratings-list">
                ${ratingsHTML}
            </div>
        `;
    } catch (error) {
        console.error('加载每日评分记录失败:', error);
        const container = document.getElementById(`daily-ratings-${goalId}`);
        container.innerHTML = '<p>加载每日评分记录失败</p>';
    }
}

// 加载并显示评论
// @param {number} goalId - 目标ID
async function loadComments(goalId) {
    try {
        // 修改这里：解构获取comments和total
        const { comments, total } = await goalAPI.getComments(goalId);
        const container = document.getElementById('comments-container');
        
        if (!comments || comments.length === 0) {
            container.innerHTML = '<p class="no-comments">暂无评论</p>';
            return;
        }
        
        // 递归渲染评论
        const renderComment = (comment, level = 0) => {
            // 限制最大层级以防止样式问题
            const currentLevel = Math.min(level, 5);
            
            // 生成回复HTML
            const repliesHTML = comment.children && comment.children.length > 0 
                ? `<div class="replies-container" id="replies-${comment.id}">
                    ${comment.children.map(child => renderComment(child, level + 1)).join('')}
                  </div>`
                : '';
            
            // 确定CSS类
            const cssClass = level === 0 ? 'comment-item' : `comment-item reply-item reply-level-${currentLevel}`;
            
            // 检查是否有子评论来决定是否显示收起按钮
            const hasReplies = comment.children && comment.children.length > 0;
            
            return `
                <div class="${cssClass}" data-comment-id="${comment.id}" id="comment-${comment.id}">
                    <div class="comment-content">${comment.content}</div>
                    <div class="comment-meta">
                        <span class="comment-date">${new Date(comment.created_at).toLocaleString()}</span>
                        ${hasReplies ? `<button class="btn btn-small toggle-reply-btn" data-comment-id="${comment.id}">收起回复</button>` : ''}
                        <button class="btn btn-small reply-btn" data-parent-id="${comment.id}">回复</button>
                    </div>
                    <!-- 回复表单（默认隐藏） -->
                    <div class="reply-form" id="reply-form-${comment.id}" style="display: none;">
                        <textarea placeholder="请输入回复内容..."></textarea>
                        <div class="reply-form-actions">
                            <button class="btn btn-small submit-reply" data-parent-id="${comment.id}">提交回复</button>
                            <button class="btn btn-small cancel-reply" data-parent-id="${comment.id}">取消</button>
                        </div>
                    </div>
                    <!-- 回复列表 -->
                    ${repliesHTML}
                </div>
            `;
        };
        
        // 生成所有评论HTML
        const commentsHTML = comments.map(comment => renderComment(comment)).join('');
        
        // 修改这里：直接使用后端返回的total字段
        container.innerHTML = `
            <h4>全部评论 (${total})</h4>
            <div class="comments-list">
                ${commentsHTML}
            </div>
        `;
        
        // 绑定回复按钮事件
        bindReplyEvents(goalId);
    } catch (error) {
        console.error('加载评论失败:', error);
        const container = document.getElementById('comments-container');
        container.innerHTML = '<p>加载评论失败</p>';
    }
}

// 绑定回复按钮事件
function bindReplyEvents(goalId) {
    // 绑定回复按钮点击事件
    document.querySelectorAll('.reply-btn').forEach(button => {
        button.addEventListener('click', function() {
            const parentId = this.getAttribute('data-parent-id');
            const replyForm = document.getElementById(`reply-form-${parentId}`);
            // 切换回复表单的显示状态
            replyForm.style.display = replyForm.style.display === 'none' ? 'block' : 'none';
            
            // 绑定提交回复按钮事件
            const submitReplyBtn = replyForm.querySelector('.submit-reply');
            if (submitReplyBtn) {
                submitReplyBtn.onclick = function() {
                    submitReply(goalId, parentId, replyForm.querySelector('textarea').value);
                };
            }
        });
    });
    
    // 绑定取消回复按钮事件
    document.querySelectorAll('.cancel-reply').forEach(button => {
        button.addEventListener('click', function() {
            const parentId = this.getAttribute('data-parent-id');
            const replyForm = document.getElementById(`reply-form-${parentId}`);
            // 隐藏回复表单
            replyForm.style.display = 'none';
            // 清空回复内容
            replyForm.querySelector('textarea').value = '';
        });
    });
    
    // 重新绑定收起/展开回复按钮事件（在重新加载评论后需要重新绑定）
    bindToggleReplyEvents();
}

// 在loadComments函数中移除bindReplyEvents的内部定义，只保留调用
async function loadComments(goalId) {
    try {
        // 修改这里：解构获取comments和total
        const { comments, total } = await goalAPI.getComments(goalId);
        const container = document.getElementById('comments-container');
        
        if (!comments || comments.length === 0) {
            container.innerHTML = '<p class="no-comments">暂无评论</p>';
            return;
        }
        
        // 递归渲染评论
        const renderComment = (comment, level = 0) => {
            // 限制最大层级以防止样式问题
            const currentLevel = Math.min(level, 5);
            
            // 生成回复HTML
            const repliesHTML = comment.children && comment.children.length > 0 
                ? `<div class="replies-container" id="replies-${comment.id}">
                    ${comment.children.map(child => renderComment(child, level + 1)).join('')}
                  </div>`
                : '';
            
            // 确定CSS类
            const cssClass = level === 0 ? 'comment-item' : `comment-item reply-item reply-level-${currentLevel}`;
            
            // 检查是否有子评论来决定是否显示收起按钮
            const hasReplies = comment.children && comment.children.length > 0;
            
            return `
                <div class="${cssClass}" data-comment-id="${comment.id}" id="comment-${comment.id}">
                    <div class="comment-content">${comment.content}</div>
                    <div class="comment-meta">
                        <span class="comment-date">${new Date(comment.created_at).toLocaleString()}</span>
                        ${hasReplies ? `<button class="btn btn-small toggle-reply-btn" data-comment-id="${comment.id}">收起回复</button>` : ''}
                        <button class="btn btn-small reply-btn" data-parent-id="${comment.id}">回复</button>
                    </div>
                    <!-- 回复表单（默认隐藏） -->
                    <div class="reply-form" id="reply-form-${comment.id}" style="display: none;">
                        <textarea placeholder="请输入回复内容..."></textarea>
                        <button class="btn btn-small submit-reply" data-parent-id="${comment.id}">提交回复</button>
                    </div>
                    <!-- 回复列表 -->
                    ${repliesHTML}
                </div>
            `;
        };
        
        // 生成所有评论HTML
        const commentsHTML = comments.map(comment => renderComment(comment)).join('');
        
        // 修改这里：直接使用后端返回的total字段
        container.innerHTML = `
            <h4>全部评论 (${total})</h4>
            <div class="comments-list">
                ${commentsHTML}
            </div>
        `;
        
        // 绑定回复按钮事件
        bindReplyEvents(goalId);
    } catch (error) {
        console.error('加载评论失败:', error);
        const container = document.getElementById('comments-container');
        container.innerHTML = '<p>加载评论失败</p>';
    }
}

// 绑定收起/展开回复按钮事件
function bindToggleReplyEvents() {
    // 为所有收起/展开回复按钮绑定事件
    document.querySelectorAll('.toggle-reply-btn').forEach(button => {
        button.addEventListener('click', function() {
            const commentId = this.getAttribute('data-comment-id');
            const repliesContainer = document.getElementById(`replies-${commentId}`);
            
            if (repliesContainer) {
                // 切换回复列表的显示状态
                const isHidden = repliesContainer.style.display === 'none';
                repliesContainer.style.display = isHidden ? 'block' : 'none';
                
                // 更新按钮文本
                this.textContent = isHidden ? '收起回复' : '展开回复';
            }
        });
    });
}

// 计算评论总数（包括回复）
function countTotalComments(comments) {
    let count = 0;
    
    function traverse(comment) {
        count++;
        if (comment.children && comment.children.length > 0) {
            comment.children.forEach(traverse);
        }
    }
    
    comments.forEach(traverse);
    return count;
}

// 提交评论
async function submitComment(goalId) {
    const content = document.getElementById('comment-content').value.trim();
    
    if (!content) {
        alert('请输入评论内容');
        return;
    }
    
    try {
        const commentData = {
            goal_id: parseInt(goalId),
            content: content
        };
        
        // 调用API创建评论
        await goalAPI.createComment(goalId, commentData);
        
        // 清空评论输入框
        document.getElementById('comment-content').value = '';
        
        // 重新加载评论
        await loadComments(goalId);
        
        alert('评论发布成功！');
    } catch (error) {
        console.error('发布评论失败:', error);
        alert('发布评论失败，请稍后重试');
    }
}

// 提交回复
async function submitReply(goalId, parentId, content) {
    content = content.trim();
    
    if (!content) {
        alert('请输入回复内容');
        return;
    }
    
    try {
        const commentData = {
            goal_id: parseInt(goalId),
            parent_id: parseInt(parentId),
            content: content
        };
        
        // 调用API创建回复
        await goalAPI.createComment(goalId, commentData);
        
        // 重新加载评论
        await loadComments(goalId);
        
        alert('回复发布成功！');
    } catch (error) {
        console.error('发布回复失败:', error);
        alert('发布回复失败，请稍后重试');
    }
}

// 创建新目标
async function createGoal(event) {
    // 阻止表单默认提交行为
    event.preventDefault();
    
    // 获取表单元素
    const form = document.getElementById('goal-form');
    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
    const category = document.getElementById('category').value;
    
    // 简单验证
    if (!title || !category) {
        alert('请填写必填字段');
        return;
    }
    
    // 准备目标数据
    const goalData = {
        title: title,
        description: description,
        category: category,
        stars: 0  // 新创建的目标默认星数为0
    };
    
    try {
        // 调用API创建目标
        const newGoal = await goalAPI.createGoal(goalData);
        
        // 显示成功消息
        alert('目标创建成功！');
        
        // 重置表单
        form.reset();
        
        // 可选：跳转到目标列表页面
        // window.location.href = 'goal-list.html';
    } catch (error) {
        console.error('创建目标失败:', error);
        alert('创建目标失败，请稍后重试');
    }
}

// 绑定类别筛选事件
document.addEventListener('DOMContentLoaded', function() {
    const categoryFilter = document.getElementById('category-filter');
    if (categoryFilter) {
        categoryFilter.addEventListener('change', function() {
            const selectedCategory = this.value;
            if (selectedCategory === 'all') {
                loadAllGoals();
            } else {
                loadGoalsByCategory(selectedCategory);
            }
        });
    }
    
    // 监听来自其他页面的消息，更新总星数
    window.addEventListener('message', function(event) {
        if (event.data.type === 'UPDATE_TOTAL_STARS') {
            loadTotalStars();
        }
    });
});


// 删除目标
async function deleteGoal(goalId) {
    try {
        // 调用API删除目标
        await goalAPI.deleteGoal(goalId);
        
        // 显示成功消息
        alert('目标删除成功！');
        
        // 跳转到目标列表页面
        window.location.href = 'goal-list.html';
    } catch (error) {
        console.error('删除目标失败:', error);
        alert('删除目标失败，请稍后重试');
    }
}