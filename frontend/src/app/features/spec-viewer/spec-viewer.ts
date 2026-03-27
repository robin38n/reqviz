import { Component, inject, type OnInit } from "@angular/core";
import { ActivatedRoute, RouterLink } from "@angular/router";
import { SpecGraphService } from "./spec-graph.service";

@Component({
	selector: "app-spec-viewer",
	standalone: true,
	imports: [RouterLink],
	template: `
    <div class="viewer">
      @if (svc.loading()) {
        <p class="status">Loading spec...</p>
      }

      @if (svc.error(); as err) {
        <div class="error">{{ err }}</div>
      }

      @if (svc.summary(); as s) {
        <header>
          <a routerLink="/" class="back">&larr; Upload</a>
          <h1>{{ s.title }} <span class="version">v{{ s.version }}</span></h1>
          <div class="stats">
            <span class="stat">{{ svc.endpointNodes().length }} endpoints</span>
            <span class="stat">{{ svc.schemaNodes().length }} schemas</span>
            <span class="stat">{{ svc.edgeCount() }} relationships</span>
          </div>
        </header>

        <div class="columns">
          <section>
            <h2>Endpoints</h2>
            @for (ep of svc.endpointNodes(); track ep.id) {
              <div class="endpoint">
                <span class="method" [attr.data-method]="ep.method">{{ ep.method }}</span>
                <span class="path">{{ ep.path }}</span>
                @if (ep.summary) {
                  <span class="summary">{{ ep.summary }}</span>
                }
              </div>
            } @empty {
              <p class="empty">No endpoints found</p>
            }
          </section>

          <section>
            <h2>Schemas</h2>
            @for (sc of svc.schemaNodes(); track sc.id) {
              <div class="schema-card">
                <strong>{{ sc.name }}</strong>
                <span class="prop-count">{{ sc.properties.length }} props</span>
                @if (sc.properties.length > 0) {
                  <div class="props">{{ sc.properties.join(', ') }}</div>
                }
              </div>
            } @empty {
              <p class="empty">No schemas found</p>
            }
          </section>
        </div>

        @if (svc.graph(); as g) {
          @if (g.edges.length > 0) {
            <section class="edges">
              <h2>Relationships</h2>
              <table>
                <thead>
                  <tr>
                    <th>Source</th>
                    <th>Kind</th>
                    <th>Target</th>
                    <th>Label</th>
                  </tr>
                </thead>
                <tbody>
                  @for (edge of g.edges; track $index) {
                    <tr>
                      <td>{{ edge.source }}</td>
                      <td><span class="edge-kind">{{ edge.kind }}</span></td>
                      <td>{{ edge.target }}</td>
                      <td>{{ edge.label ?? '' }}</td>
                    </tr>
                  }
                </tbody>
              </table>
            </section>
          }
        }
      }
    </div>
  `,
	styles: `
    .viewer {
      max-width: 960px;
      margin: 2rem auto;
      padding: 1rem;
      font-family: system-ui, sans-serif;
    }
    .status {
      color: #666;
    }
    .error {
      padding: 0.75rem;
      background: #fef2f2;
      color: #dc2626;
      border-radius: 4px;
    }
    .back {
      color: #666;
      text-decoration: none;
      font-size: 0.875rem;
    }
    .back:hover {
      color: #333;
    }
    header h1 {
      margin: 0.5rem 0 0;
    }
    .version {
      color: #666;
      font-size: 0.875rem;
    }
    .stats {
      display: flex;
      gap: 1rem;
      margin-top: 0.5rem;
    }
    .stat {
      padding: 0.25rem 0.5rem;
      background: #f3f4f6;
      border-radius: 4px;
      font-size: 0.875rem;
      color: #555;
    }
    .columns {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1.5rem;
      margin-top: 1.5rem;
    }
    h2 {
      font-size: 1rem;
      margin: 0 0 0.75rem;
      color: #333;
    }
    .endpoint {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem;
      border: 1px solid #e5e7eb;
      border-radius: 4px;
      margin-bottom: 0.5rem;
      font-size: 0.875rem;
    }
    .method {
      font-weight: 600;
      font-size: 0.75rem;
      padding: 0.125rem 0.375rem;
      border-radius: 3px;
      color: #fff;
      background: #6b7280;
    }
    .method[data-method="GET"] { background: #16a34a; }
    .method[data-method="POST"] { background: #2563eb; }
    .method[data-method="PUT"] { background: #d97706; }
    .method[data-method="PATCH"] { background: #9333ea; }
    .method[data-method="DELETE"] { background: #dc2626; }
    .path {
      font-family: monospace;
      font-weight: 500;
    }
    .summary {
      color: #888;
      font-size: 0.8rem;
    }
    .schema-card {
      padding: 0.5rem;
      border: 1px solid #e5e7eb;
      border-radius: 4px;
      margin-bottom: 0.5rem;
      font-size: 0.875rem;
    }
    .prop-count {
      margin-left: 0.5rem;
      color: #888;
      font-size: 0.75rem;
    }
    .props {
      margin-top: 0.25rem;
      font-family: monospace;
      font-size: 0.75rem;
      color: #666;
    }
    .edges {
      margin-top: 1.5rem;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      font-size: 0.8rem;
    }
    th, td {
      text-align: left;
      padding: 0.375rem 0.5rem;
      border-bottom: 1px solid #e5e7eb;
    }
    th {
      font-weight: 600;
      color: #555;
    }
    .edge-kind {
      font-family: monospace;
      padding: 0.125rem 0.25rem;
      background: #f3f4f6;
      border-radius: 3px;
      font-size: 0.7rem;
    }
    .empty {
      color: #999;
      font-size: 0.875rem;
    }
  `,
})
export class SpecViewerComponent implements OnInit {
	protected readonly svc = inject(SpecGraphService);
	private readonly route = inject(ActivatedRoute);

	ngOnInit() {
		const id = this.route.snapshot.paramMap.get("id");
		if (id) {
			this.svc.loadSpec(id);
		}
	}
}
