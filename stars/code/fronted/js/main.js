// 主脚本文件

// 加载总星数
async function loadTotalStars() {
    try {
        const result = await goalAPI.getTotalStars();
        document.getElementById('total-stars').textContent = result.total_stars;
    } catch (error) {
        console.error('加载总星数失败:', error);
    }
}

// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    console.log('星目标管理系统已加载');
    
    // 加载总星数
    loadTotalStars();
    
    // 可以在这里添加全局的初始化代码
    // 例如：检查用户登录状态、初始化导航等
});