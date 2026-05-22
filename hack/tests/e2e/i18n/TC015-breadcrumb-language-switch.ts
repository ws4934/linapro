import { test, expect } from '../../fixtures/auth';
import { waitForRouteReady } from '../../support/ui';

test.describe('TC-15 breadcrumb language switch', () => {
  test('TC-15a: breadcrumb titles relocalize immediately after language switch', async ({
    adminPage,
    mainLayout,
  }) => {
    await mainLayout.switchLanguage('简体中文');

    await adminPage.goto('/dashboard/workspace', {
      waitUntil: 'domcontentloaded',
    });
    await waitForRouteReady(adminPage);

    await expect(mainLayout.breadcrumbItem('工作台')).toBeVisible();

    await adminPage.goto('/system/user', { waitUntil: 'domcontentloaded' });
    await waitForRouteReady(adminPage);

    await expect(mainLayout.breadcrumbItem('权限管理')).toBeVisible();
    await expect(mainLayout.breadcrumbItem('用户管理')).toBeVisible();

    await mainLayout.switchLanguage('English');

    await expect(mainLayout.breadcrumbItem('Access')).toBeVisible();
    await expect(mainLayout.breadcrumbItem('Users')).toBeVisible();
    await expect(mainLayout.breadcrumbItem('权限管理')).toHaveCount(0);
    await expect(mainLayout.breadcrumbItem('用户管理')).toHaveCount(0);

    await mainLayout.switchLanguage('简体中文');

    await expect(mainLayout.breadcrumbItem('权限管理')).toBeVisible();
    await expect(mainLayout.breadcrumbItem('用户管理')).toBeVisible();
    await expect(mainLayout.breadcrumbItem('Access')).toHaveCount(0);
    await expect(mainLayout.breadcrumbItem('Users')).toHaveCount(0);
  });
});
