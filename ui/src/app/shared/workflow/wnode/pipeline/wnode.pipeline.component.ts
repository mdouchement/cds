import {Component, Input} from '@angular/core';
import {PipelineStatus} from '../../../../model/pipeline.model';
import {WNode, Workflow} from '../../../../model/workflow.model';
import {WorkflowNodeRun} from '../../../../model/workflow.run.model';
import {AutoUnsubscribe} from '../../../decorator/autoUnsubscribe';

@Component({
    selector: 'app-workflow-wnode-pipeline',
    templateUrl: './node.pipeline.html',
    styleUrls: ['./node.pipeline.scss']
})
@AutoUnsubscribe()
export class WorkflowWNodePipelineComponent {

    @Input() public node: WNode;
    @Input() public workflow: Workflow;
    @Input() public noderun: WorkflowNodeRun;
    @Input() public selected: boolean;

    pipelineStatus = PipelineStatus;
}