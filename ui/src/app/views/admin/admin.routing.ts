import { ModuleWithProviders } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { CanActivateAuthAdminRoute } from '../../service/auth/authenAdminRouteActivate';
import { AdminComponent } from './admin.component';
import { BroadcastAddComponent } from './broadcast/add/broadcast.add.component';
import { BroadcastEditComponent } from './broadcast/edit/broadcast.edit.component';
import { BroadcastListComponent } from './broadcast/list/broadcast.list.component';
import { HooksTaskComponent } from './hooks-tasks/details/hooks-task.component';
import { HooksTasksComponent } from './hooks-tasks/hooks-tasks.component';
import { ServiceComponent } from './services/service/service.component';
import { ServicesComponent } from './services/services.component';
import { WorkerModelPatternAddComponent } from './worker-model-pattern/add/worker-model-pattern.add.component';
import { WorkerModelPatternEditComponent } from './worker-model-pattern/edit/worker-model-pattern.edit.component';
import { WorkerModelPatternComponent } from './worker-model-pattern/worker-model-pattern.component';

const routes: Routes = [
    {
        path: '',
        component: AdminComponent,
        canActivateChild: [CanActivateAuthAdminRoute],
        canActivate: [CanActivateAuthAdminRoute],
        children: [
            {
                path: 'worker-model-pattern',
                component: WorkerModelPatternComponent,
                data: { title: 'List • Worker Model Pattern' }
            },
            {
                path: 'worker-model-pattern/add',
                component: WorkerModelPatternAddComponent,
                data: { title: 'Add • Worker Model Pattern' }
            },
            {
                path: 'worker-model-pattern/:type/:name',
                component: WorkerModelPatternEditComponent,
                data: { title: '{name} • Edit • Worker Model Pattern' }
            },
            {
                path: 'broadcast',
                component: BroadcastListComponent,
                data: { title: 'List • Broadcast' }
            },
            {
                path: 'broadcast/add',
                component: BroadcastAddComponent,
                data: { title: 'Add • Broadcast' }
            },
            {
                path: 'broadcast/:id',
                component: BroadcastEditComponent,
                data: { title: 'Edit {id} • Broadcast' }
            },
            {
                path: 'hooks-tasks',
                component: HooksTasksComponent,
                data: { title: 'Hooks tasks summary' }
            },
            {
                path: 'hooks-tasks/:id',
                component: HooksTaskComponent,
                data: { title: 'Hooks task details' }
            },
            {
                path: 'services',
                component: ServicesComponent,
                data: { title: 'Services' }
            },
            {
                path: 'services/:name',
                component: ServiceComponent,
                data: { title: 'Service' }
            }
        ]
    }
];

export const AdminRouting: ModuleWithProviders = RouterModule.forChild(routes);
