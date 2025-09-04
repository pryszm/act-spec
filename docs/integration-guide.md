# ASTRA Integration Guide

This guide provides comprehensive patterns and practices for integrating ASTRA into production applications and business systems. Whether you're building a new conversational application or adding ASTRA support to existing systems, these patterns will help you achieve reliable, scalable, and maintainable integrations.

## Integration Architecture Overview

ASTRA acts as the universal interchange format for conversational state, enabling different systems to interoperate through structured, type-safe interfaces. The integration architecture supports multiple patterns depending on your scale, reliability, and performance requirements.

### Key Integration Principles

**Event-Driven Architecture** - ASTRA acts flow through event streams, enabling loose coupling and horizontal scaling

**Idempotent Operations** - All integrations support retry and recovery patterns without side effects

**Observability by Design** - Complete traceability from conversation to business system execution

**Gradual Adoption** - Incremental integration paths that don't require big-bang migrations

## Integration Patterns

### 1. Event-Driven Integration

The most common and scalable pattern for ASTRA integration uses event streams to decouple conversation processing from business system integration.

**Architecture:**
```
Conversation -> ASTRA Acts -> Event Stream -> Business Systems
              ↓
         Act Validation
         State Management
         Routing Logic
```

**Implementation Example:**

```typescript
import { ConversationAct, isCommit, isFact } from '@astra/model-ts';

// Event publisher
class ConversationEventPublisher {
  constructor(
    private eventBus: EventBus,
    private actValidator: ActValidator
  ) {}

  async publishAct(act: ConversationAct): Promise<void> {
    // Validate act structure
    const validation = await this.actValidator.validate(act);
    if (!validation.valid) {
      throw new ValidationError('Invalid act structure', validation.errors);
    }

    // Publish to appropriate topic based on act type
    const topic = this.getTopicForAct(act);
    await this.eventBus.publish(topic, {
      conversationId: act.metadata?.conversation_id,
      act,
      timestamp: act.timestamp,
      traceId: generateTraceId()
    });
  }

  private getTopicForAct(act: ConversationAct): string {
    return `conversation.${act.type}`;
  }
}

// Business system integration subscriber
class OrderManagementSubscriber {
  constructor(
    private orderService: OrderService,
    private eventBus: EventBus
  ) {
    this.setupSubscriptions();
  }

  private setupSubscriptions(): void {
    // Handle order-related facts
    this.eventBus.subscribe('conversation.fact', this.handleFact.bind(this));
    
    // Handle order commits
    this.eventBus.subscribe('conversation.commit', this.handleCommit.bind(this));
  }

  private async handleFact(event: ConversationEvent): Promise<void> {
    if (!isFact(event.act)) return;

    const fact = event.act;
    
    // Only process order-related entities
    if (typeof fact.entity === 'object' && fact.entity.type === 'order') {
      await this.updateOrderEntity(fact);
    }
  }

  private async handleCommit(event: ConversationEvent): Promise<void> {
    if (!isCommit(event.act)) return;

    const commit = event.act;
    
    // Only process commits targeting order management
    if (commit.system === 'order_management') {
      await this.executeOrderAction(commit);
    }
  }

  private async updateOrderEntity(fact: Fact): Promise<void> {
    try {
      const orderId = typeof fact.entity === 'string' 
        ? fact.entity 
        : fact.entity.id;

      await this.orderService.updateField(
        orderId, 
        fact.field, 
        fact.value
      );

      // Publish success event
      await this.eventBus.publish('order.field.updated', {
        orderId,
        field: fact.field,
        value: fact.value,
        traceId: event.traceId
      });
    } catch (error) {
      // Publish error for handling
      await this.eventBus.publish('order.field.error', {
        error: error.message,
        fact,
        traceId: event.traceId
      });
    }
  }

  private async executeOrderAction(commit: Commit): Promise<void> {
    try {
      const orderId = typeof commit.entity === 'string'
        ? commit.entity
        : commit.entity.id;

      let result;
      switch (commit.action) {
        case 'create':
          result = await this.orderService.createOrder(orderId);
          break;
        case 'update':
          result = await this.orderService.updateOrder(orderId);
          break;
        case 'delete':
          result = await this.orderService.cancelOrder(orderId);
          break;
        default:
          throw new Error(`Unsupported action: ${commit.action}`);
      }

      // Publish success with transaction ID
      await this.eventBus.publish('order.commit.success', {
        orderId,
        action: commit.action,
        transactionId: result.transactionId,
        traceId: event.traceId
      });
    } catch (error) {
      // Publish error for retry handling
      await this.eventBus.publish('order.commit.error', {
        error: error.message,
        commit,
        recoverable: this.isRecoverableError(error),
        traceId: event.traceId
      });
    }
  }
}
```

### 2. Synchronous API Integration

For applications requiring immediate feedback or simpler deployment models, synchronous API integration provides direct request-response patterns.

```typescript
// REST API integration
class ConversationAPI {
  constructor(
    private conversationProcessor: ConversationProcessor,
    private businessIntegrations: BusinessIntegrationRegistry
  ) {}

  async processAct(req: Request, res: Response): Promise<void> {
    try {
      const act = req.body as ConversationAct;
      
      // Validate act
      const validation = await this.conversationProcessor.validate(act);
      if (!validation.valid) {
        return res.status(400).json({
          error: 'Invalid act',
          details: validation.errors
        });
      }

      // Process act and trigger integrations
      const result = await this.conversationProcessor.process(act);
      
      // Execute any required business system calls
      const integrationResults = await this.executeIntegrations(act, result);

      res.json({
        processedAct: result.act,
        integrations: integrationResults,
        nextActions: result.nextActions
      });
    } catch (error) {
      res.status(500).json({
        error: 'Processing failed',
        message: error.message
      });
    }
  }

  private async executeIntegrations(
    act: ConversationAct, 
    processResult: ProcessResult
  ): Promise<IntegrationResult[]> {
    const results: IntegrationResult[] = [];

    // Execute integrations based on act type and content
    for (const integration of this.getRequiredIntegrations(act)) {
      try {
        const result = await integration.execute(act, processResult);
        results.push({
          system: integration.name,
          status: 'success',
          result
        });
      } catch (error) {
        results.push({
          system: integration.name,
          status: 'error',
          error: error.message,
          recoverable: integration.isRecoverable(error)
        });
      }
    }

    return results;
  }
}
```

### 3. Stream Processing Integration

For high-throughput applications or complex event processing requirements, stream processing enables sophisticated conversation analysis and real-time business intelligence.

```python
# Apache Kafka + Python stream processing example
from kafka import KafkaConsumer, KafkaProducer
from astra_model import ConversationAct, Fact, Commit
import json
import logging

class ConversationStreamProcessor:
    def __init__(self, kafka_config):
        self.consumer = KafkaConsumer(
            'conversation-acts',
            bootstrap_servers=kafka_config['bootstrap_servers'],
            value_deserializer=lambda m: json.loads(m.decode('utf-8')),
            group_id='conversation-processor'
        )
        
        self.producer = KafkaProducer(
            bootstrap_servers=kafka_config['bootstrap_servers'],
            value_serializer=lambda v: json.dumps(v).encode('utf-8')
        )
        
        self.entity_state_store = EntityStateStore()
        self.business_rules_engine = BusinessRulesEngine()
        
    def start_processing(self):
        """Start the stream processing loop"""
        for message in self.consumer:
            try:
                self.process_message(message.value)
            except Exception as e:
                logging.error(f"Error processing message: {e}")
                self.handle_processing_error(message.value, e)
    
    def process_message(self, event_data):
        """Process a single conversation act event"""
        act_data = event_data['act']
        conversation_id = event_data.get('conversationId')
        trace_id = event_data.get('traceId')
        
        # Deserialize act based on type
        act = self.deserialize_act(act_data)
        
        # Update entity state
        if isinstance(act, Fact):
            self.update_entity_state(act, conversation_id)
        
        # Execute business rules
        rule_results = self.business_rules_engine.evaluate(act, conversation_id)
        
        # Trigger downstream actions
        if isinstance(act, Commit):
            self.handle_commit_act(act, conversation_id, trace_id)
        
        # Publish derived events
        for derived_event in rule_results.derived_events:
            self.producer.send(
                derived_event['topic'], 
                derived_event['payload']
            )
    
    def update_entity_state(self, fact: Fact, conversation_id: str):
        """Update entity state based on fact"""
        entity_id = fact.entity if isinstance(fact.entity, str) else fact.entity.id
        
        # Get current state
        current_state = self.entity_state_store.get(entity_id)
        
        # Apply fact operation
        updated_state = self.apply_fact_operation(current_state, fact)
        
        # Validate against business constraints
        validation_result = self.business_rules_engine.validate_entity_state(
            entity_id, 
            updated_state
        )
        
        if validation_result.valid:
            # Store updated state
            self.entity_state_store.put(entity_id, updated_state)
            
            # Publish state change event
            self.producer.send('entity-state-changed', {
                'entityId': entity_id,
                'conversationId': conversation_id,
                'field': fact.field,
                'previousValue': fact.previous_value,
                'newValue': fact.value,
                'timestamp': fact.timestamp
            })
        else:
            # Publish validation error
            self.producer.send('entity-validation-error', {
                'entityId': entity_id,
                'conversationId': conversation_id,
                'errors': validation_result.errors,
                'fact': fact.model_dump()
            })
    
    def handle_commit_act(self, commit: Commit, conversation_id: str, trace_id: str):
        """Handle commit acts by triggering business system integration"""
        entity_id = commit.entity if isinstance(commit.entity, str) else commit.entity.id
        
        # Get final entity state
        entity_state = self.entity_state_store.get(entity_id)
        
        # Prepare integration payload
        integration_payload = {
            'commitId': commit.id,
            'conversationId': conversation_id,
            'traceId': trace_id,
            'entity': entity_state,
            'action': commit.action,
            'system': commit.system,
            'idempotencyKey': commit.idempotency_key
        }
        
        # Route to appropriate integration topic
        integration_topic = f"integration-{commit.system}" if commit.system else "integration-default"
        
        self.producer.send(integration_topic, integration_payload)
```

## Business System Integration Strategies

### Customer Relationship Management (CRM)

```typescript
class CRMIntegration {
  constructor(private crmClient: CRMClient) {}

  async handleCustomerFacts(facts: Fact[]): Promise<void> {
    // Group facts by customer entity
    const customerFacts = facts.filter(f => 
      typeof f.entity === 'object' && f.entity.type === 'customer'
    );

    for (const fact of customerFacts) {
      await this.updateCustomerField(fact);
    }
  }

  private async updateCustomerField(fact: Fact): Promise<void> {
    const customerId = typeof fact.entity === 'string' 
      ? fact.entity 
      : fact.entity.external_id || fact.entity.id;

    const crmFieldMapping = {
      'email': 'email_address',
      'phone': 'primary_phone',
      'company': 'company_name',
      'address': 'billing_address'
    };

    const crmField = crmFieldMapping[fact.field] || fact.field;

    try {
      await this.crmClient.updateContact(customerId, {
        [crmField]: fact.value
      });

      // Record successful update
      await this.recordIntegrationSuccess('crm', 'update_contact', {
        customerId,
        field: crmField,
        value: fact.value
      });
    } catch (error) {
      await this.recordIntegrationError('crm', 'update_contact', error, {
        customerId,
        field: crmField,
        fact
      });
    }
  }

  async handleCustomerCommits(commit: Commit): Promise<void> {
    if (commit.system !== 'crm') return;

    const customerId = typeof commit.entity === 'string'
      ? commit.entity
      : commit.entity.external_id || commit.entity.id;

    try {
      switch (commit.action) {
        case 'create':
          await this.createCustomerRecord(customerId, commit);
          break;
        case 'update':
          await this.updateCustomerRecord(customerId, commit);
          break;
        default:
          throw new Error(`Unsupported CRM action: ${commit.action}`);
      }
    } catch (error) {
      await this.handleCommitError(commit, error);
    }
  }
}
```

### Order Management Systems

```python
class OrderManagementIntegration:
    def __init__(self, order_service, inventory_service, payment_service):
        self.order_service = order_service
        self.inventory_service = inventory_service
        self.payment_service = payment_service
        
    async def handle_order_commit(self, commit: Commit, entity_state: dict):
        """Handle order-related commits with full business logic"""
        order_id = commit.entity if isinstance(commit.entity, str) else commit.entity.id
        
        try:
            if commit.action == 'create':
                await self._create_order_workflow(order_id, entity_state, commit)
            elif commit.action == 'update':
                await self._update_order_workflow(order_id, entity_state, commit)
            elif commit.action == 'cancel':
                await self._cancel_order_workflow(order_id, entity_state, commit)
                
        except Exception as e:
            await self._handle_order_error(order_id, commit, e)
    
    async def _create_order_workflow(self, order_id: str, state: dict, commit: Commit):
        """Complete order creation workflow with validation and inventory checks"""
        
        # 1. Validate order completeness
        required_fields = ['customer_id', 'items', 'total_amount', 'delivery_address']
        missing_fields = [f for f in required_fields if f not in state]
        
        if missing_fields:
            raise ValidationError(f"Missing required fields: {missing_fields}")
        
        # 2. Check inventory availability
        for item in state['items']:
            available = await self.inventory_service.check_availability(
                item['product_id'], 
                item['quantity']
            )
            if not available:
                raise InventoryError(f"Insufficient inventory for {item['product_id']}")
        
        # 3. Reserve inventory
        reservation_ids = []
        try:
            for item in state['items']:
                reservation_id = await self.inventory_service.reserve_inventory(
                    item['product_id'], 
                    item['quantity'],
                    order_id
                )
                reservation_ids.append(reservation_id)
            
            # 4. Create order record
            order_record = await self.order_service.create_order({
                'id': order_id,
                'customer_id': state['customer_id'],
                'items': state['items'],
                'total_amount': state['total_amount'],
                'delivery_address': state['delivery_address'],
                'status': 'confirmed',
                'inventory_reservations': reservation_ids
            })
            
            # 5. Process payment if payment method provided
            if 'payment_method' in state:
                await self._process_order_payment(order_record, state['payment_method'])
            
            # 6. Update commit with success
            await self._update_commit_status(commit.id, 'success', {
                'transaction_id': order_record.id,
                'reservations': reservation_ids
            })
            
        except Exception as e:
            # Rollback reservations on failure
            for reservation_id in reservation_ids:
                await self.inventory_service.cancel_reservation(reservation_id)
            raise e
    
    async def _process_order_payment(self, order: dict, payment_method: dict):
        """Process payment for the order"""
        payment_result = await self.payment_service.charge(
            amount=order['total_amount'],
            currency=order.get('currency', 'USD'),
            payment_method=payment_method,
            order_id=order['id']
        )
        
        if payment_result.status == 'failed':
            raise PaymentError(f"Payment failed: {payment_result.error_message}")
        
        # Update order with payment information
        await self.order_service.update_order(order['id'], {
            'payment_status': 'paid',
            'payment_id': payment_result.transaction_id
        })
```

### Analytics and Business Intelligence

```go
// Go example for high-performance analytics integration
package analytics

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/pryszm/astra-model-go"
)

type AnalyticsIntegration struct {
    clickhouse    ClickHouseClient
    redis         RedisClient
    timeSeriesDB  TimeSeriesDB
}

type ConversationMetrics struct {
    ConversationID string            `json:"conversation_id"`
    Timestamp     time.Time         `json:"timestamp"`
    ActType       string            `json:"act_type"`
    Speaker       string            `json:"speaker"`
    EntityType    string            `json:"entity_type,omitempty"`
    Duration      int64             `json:"duration_ms"`
    Metadata      map[string]interface{} `json:"metadata"`
}

func (a *AnalyticsIntegration) ProcessActForAnalytics(
    ctx context.Context, 
    act astra.ConversationAct,
    conversationID string,
) error {
    // Extract metrics from act
    metrics := a.extractMetrics(act, conversationID)
    
    // Real-time metrics to Redis
    if err := a.updateRealTimeMetrics(ctx, metrics); err != nil {
        return fmt.Errorf("failed to update real-time metrics: %w", err)
    }
    
    // Batch analytics to ClickHouse
    if err := a.storeAnalyticsEvent(ctx, metrics); err != nil {
        return fmt.Errorf("failed to store analytics event: %w", err)
    }
    
    // Time series data for dashboards
    if err := a.updateTimeSeries(ctx, metrics); err != nil {
        return fmt.Errorf("failed to update time series: %w", err)
    }
    
    return nil
}

func (a *AnalyticsIntegration) extractMetrics(
    act astra.ConversationAct, 
    conversationID string,
) ConversationMetrics {
    baseAct := act.GetAct()
    
    metrics := ConversationMetrics{
        ConversationID: conversationID,
        Timestamp:     baseAct.Timestamp,
        ActType:       string(baseAct.Type),
        Speaker:       baseAct.Speaker,
        Metadata:      make(map[string]interface{}),
    }
    
    // Add act-specific metrics
    switch a := act.(type) {
    case astra.Ask:
        metrics.Metadata["field"] = a.Field
        metrics.Metadata["required"] = a.Required != nil && *a.Required
        metrics.Metadata["constraints_count"] = len(a.Constraints)
        
    case astra.Fact:
        if entity, ok := a.Entity.(astra.Entity); ok {
            metrics.EntityType = entity.Type
        }
        metrics.Metadata["field"] = a.Field
        metrics.Metadata["operation"] = string(*a.Operation)
        
    case astra.Commit:
        if entity, ok := a.Entity.(astra.Entity); ok {
            metrics.EntityType = entity.Type
        }
        metrics.Metadata["action"] = string(a.Action)
        metrics.Metadata["system"] = a.System
        metrics.Metadata["status"] = string(*a.Status)
        
    case astra.Error:
        metrics.Metadata["error_code"] = a.Code
        metrics.Metadata["recoverable"] = a.Recoverable
        metrics.Metadata["severity"] = string(*a.Severity)
    }
    
    return metrics
}

func (a *AnalyticsIntegration) updateRealTimeMetrics(
    ctx context.Context, 
    metrics ConversationMetrics,
) error {
    // Update conversation act counters
    pipe := a.redis.Pipeline()
    
    // Increment act type counters
    pipe.Incr(ctx, fmt.Sprintf("acts:%s:count", metrics.ActType))
    pipe.Incr(ctx, fmt.Sprintf("conversations:%s:acts", metrics.ConversationID))
    
    // Update hourly metrics
    hourKey := fmt.Sprintf("acts:hourly:%s", 
        metrics.Timestamp.Format("2006010215"))
    pipe.Incr(ctx, hourKey)
    pipe.Expire(ctx, hourKey, 48*time.Hour)
    
    // Entity type metrics
    if metrics.EntityType != "" {
        pipe.Incr(ctx, fmt.Sprintf("entities:%s:count", metrics.EntityType))
    }
    
    _, err := pipe.Exec(ctx)
    return err
}

func (a *AnalyticsIntegration) storeAnalyticsEvent(
    ctx context.Context,
    metrics ConversationMetrics,
) error {
    // Prepare data for ClickHouse
    data := map[string]interface{}{
        "conversation_id": metrics.ConversationID,
        "timestamp":      metrics.Timestamp,
        "act_type":       metrics.ActType,
        "speaker":        metrics.Speaker,
        "entity_type":    metrics.EntityType,
        "metadata":       metrics.Metadata,
    }
    
    // Insert into ClickHouse for analytical queries
    query := `
        INSERT INTO conversation_acts 
        (conversation_id, timestamp, act_type, speaker, entity_type, metadata)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    
    metadataJSON, _ := json.Marshal(metrics.Metadata)
    
    return a.clickhouse.Exec(ctx, query,
        metrics.ConversationID,
        metrics.Timestamp,
        metrics.ActType,
        metrics.Speaker,
        metrics.EntityType,
        string(metadataJSON),
    )
}
```

## Production Deployment Patterns

### Microservices Architecture

```yaml
# docker-compose.yml for microservices deployment
version: '3.8'

services:
  # Core conversation processing
  conversation-processor:
    image: myorg/conversation-processor:latest
    environment:
      - KAFKA_BROKERS=kafka:9092
      - REDIS_URL=redis:6379
      - DATABASE_URL=postgres://postgres:password@postgres:5432/conversations
    depends_on:
      - kafka
      - redis
      - postgres

  # Act validation service
  act-validator:
    image: myorg/act-validator:latest
    environment:
      - SCHEMA_REGISTRY_URL=http://schema-registry:8081
    depends_on:
      - schema-registry

  # Business system integrations
  crm-integration:
    image: myorg/crm-integration:latest
    environment:
      - KAFKA_BROKERS=kafka:9092
      - CRM_API_URL=https://api.crm-system.com
      - CRM_API_KEY_SECRET=crm-api-key
    depends_on:
      - kafka
    secrets:
      - crm-api-key

  order-integration:
    image: myorg/order-integration:latest
    environment:
      - KAFKA_BROKERS=kafka:9092
      - ORDER_DB_URL=postgres://postgres:password@order-db:5432/orders
    depends_on:
      - kafka
      - order-db

  # Analytics pipeline
  analytics-processor:
    image: myorg/analytics-processor:latest
    environment:
      - KAFKA_BROKERS=kafka:9092
      - CLICKHOUSE_URL=http://clickhouse:8123
      - REDIS_URL=redis:6379
    depends_on:
      - kafka
      - clickhouse
      - redis

  # Infrastructure
  kafka:
    image: confluentinc/cp-kafka:latest
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
    depends_on:
      - zookeeper

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  schema-registry:
    image: confluentinc/cp-schema-registry:latest
    environment:
      SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS: kafka:9092

  redis:
    image: redis:alpine

  postgres:
    image: postgres:13
    environment:
      POSTGRES_DB: conversations
      POSTGRES_PASSWORD: password

  clickhouse:
    image: yandex/clickhouse-server:latest

secrets:
  crm-api-key:
    external: true
```

### Kubernetes Deployment

```yaml
# kubernetes/conversation-processor-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: conversation-processor
  labels:
    app: conversation-processor
spec:
  replicas: 3
  selector:
    matchLabels:
      app: conversation-processor
  template:
    metadata:
      labels:
        app: conversation-processor
    spec:
      containers:
      - name: conversation-processor
        image: myorg/conversation-processor:v1.2.0
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
        env:
        - name: KAFKA_BROKERS
          value: "kafka-cluster:9092"
        - name: REDIS_URL
          value: "redis://redis-cluster:6379"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: postgres-credentials
              key: database-url
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        volumeMounts:
        - name: config
          mountPath: /app/config
          readOnly: true
      volumes:
      - name: config
        configMap:
          name: conversation-processor-config
---
apiVersion: v1
kind: Service
metadata:
  name: conversation-processor-service
spec:
  selector:
    app: conversation-processor
  ports:
  - name: http
    port: 80
    targetPort: 8080
  - name: metrics
    port: 9090
    targetPort: 9090
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: conversation-processor-config
data:
  config.yaml: |
    server:
      port: 8080
      metrics_port: 9090
    kafka:
      brokers: ["kafka-cluster:9092"]
      consumer_group: "conversation-processor"
    validation:
      strict_mode: true
      schema_cache_ttl: "5m"
    integrations:
      timeout: "30s"
      max_retries: 3
      retry_backoff: "1s"
```

## Error Handling and Resilience

### Retry and Circuit Breaker Patterns

```typescript
import { CircuitBreaker } from 'opossum';

class ResilientIntegration {
  private circuitBreakers: Map<string, CircuitBreaker>;

  constructor() {
    this.circuitBreakers = new Map();
    this.setupCircuitBreakers();
  }

  private setupCircuitBreakers(): void {
    const circuitBreakerOptions = {
      timeout: 30000,           // 30 second timeout
      errorThresholdPercentage: 50, // Open circuit at 50% error rate
      resetTimeout: 60000,      // Try again after 1 minute
      rollingCountTimeout: 10000, // 10 second rolling window
      rollingCountBuckets: 10   // 10 buckets in the window
    };

    // Create circuit breakers for each integration
    this.circuitBreakers.set('crm', 
      new CircuitBreaker(this.callCRMAPI.bind(this), circuitBreakerOptions));
    this.circuitBreakers.set('order-management', 
      new CircuitBreaker(this.callOrderAPI.bind(this), circuitBreakerOptions));
  }

  async executeWithResilience(
    system: string, 
    operation: () => Promise<any>,
    maxRetries = 3,
    backoffMs = 1000
  ): Promise<any> {
    const circuitBreaker = this.circuitBreakers.get(system);
    if (!circuitBreaker) {
      throw new Error(`No circuit breaker configured for system: ${system}`);
    }

    let lastError;
    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        return await circuitBreaker.fire(operation);
      } catch (error) {
        lastError = error;
        
        // Don't retry if circuit is open
        if (circuitBreaker.opened) {
          throw new CircuitOpenError(`Circuit breaker open for ${system}`);
        }

        // Don't retry on certain error types
        if (!this.isRetryableError(error)) {
          throw error;
        }

        if (attempt < maxRetries) {
          const delay = backoffMs * Math.pow(2, attempt); // Exponential backoff
          await this.sleep(delay);
        }
      }
    }

    throw new MaxRetriesExceededError(`Max retries exceeded for ${system}`, lastError);
  }

  private isRetryableError(error: any): boolean {
    // Don't retry client errors (4xx) or authentication errors
    if (error.status && error.status >= 400 && error.status < 500) {
      return false;
    }
    
    // Don't retry validation errors
    if (error instanceof ValidationError) {
      return false;
    }

    // Retry on network errors, timeouts, and server errors (5xx)
    return true;
  }

  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}
```

### Dead Letter Queue Pattern

```python
class DeadLetterQueueHandler:
    """Handles failed message processing with DLQ pattern"""
    
    def __init__(self, main_queue, dlq_queue, max_retries=3):
        self.main_queue = main_queue
        self.dlq_queue = dlq_queue
        self.max_retries = max_retries
        
    async def process_with_dlq(self, message_handler, message):
        """Process message with automatic DLQ routing on failure"""
        retry_count = message.get('retryCount', 0)
        
        try:
            await message_handler(message)
            
        except Exception as e:
            if retry_count >= self.max_retries:
                # Send to dead letter queue
                await self.send_to_dlq(message, e)
            else:
                # Retry with exponential backoff
                await self.retry_message(message, retry_count + 1)
                
    async def send_to_dlq(self, message, error):
        """Send failed message to dead letter queue"""
        dlq_message = {
            'originalMessage': message,
            'error': str(error),
            'errorType': type(error).__name__,
            'timestamp': datetime.utcnow().isoformat(),
            'retryCount': message.get('retryCount', 0)
        }
        
        await self.dlq_queue.send(dlq_message)
        
        # Alert operations team
        await self.send_alert(f"Message sent to DLQ: {error}")
        
    async def retry_message(self, message, retry_count):
        """Retry message with exponential backoff"""
        message['retryCount'] = retry_count
        delay_seconds = min(300, 2 ** retry_count)  # Cap at 5 minutes
        
        await self.main_queue.send_delayed(message, delay_seconds)
```

## Monitoring and Observability

### Distributed Tracing

```typescript
import { trace, context, SpanStatusCode } from '@opentelemetry/api';

class TracedConversationProcessor {
  private tracer = trace.getTracer('conversation-processor');

  async processAct(act: ConversationAct, conversationId: string): Promise<void> {
    const span = this.tracer.startSpan('process-act', {
      attributes: {
        'act.type': act.type,
        'act.id': act.id,
        'conversation.id': conversationId,
        'speaker.id': act.speaker
      }
    });

    return context.with(trace.setSpan(context.active(), span), async () => {
      try {
        // Validate act
        await this.validateAct(act);
        
        // Process based on act type
        if (isCommit(act)) {
          await this.processCommitAct(act, conversationId);
        } else if (isFact(act)) {
          await this.processFactAct(act, conversationId);
        }

        span.setStatus({ code: SpanStatusCode.OK });
      } catch (error) {
        span.recordException(error);
        span.setStatus({ 
          code: SpanStatusCode.ERROR, 
          message: error.message 
        });
        throw error;
      } finally {
        span.end();
      }
    });
  }

  private async processCommitAct(commit: Commit, conversationId: string): Promise<void> {
    const span = this.tracer.startSpan('process-commit', {
      attributes: {
        'commit.action': commit.action,
        'commit.system': commit.system || 'unknown',
        'commit.entity': typeof commit.entity === 'string' ? commit.entity : commit.entity.id
      }
    });

    return context.with(trace.setSpan(context.active(), span), async () => {
      try {
        // Execute business system integration
        const result = await this.executeBusinessIntegration(commit);
        
        span.setAttributes({
          'commit.transaction_id': result.transactionId,
          'commit.status': 'success'
        });
      } catch (error) {
        span.recordException(error);
        span.setAttributes({
          'commit.status': 'failed',
          'commit.error': error.message
        });
        throw error;
      } finally {
        span.end();
      }
    });
  }
}
```

### Metrics and Alerting

```go
package monitoring

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // Act processing metrics
    actsProcessedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "astra_acts_processed_total",
            Help: "Total number of acts processed by type",
        },
        []string{"act_type", "status"},
    )
    
    actProcessingDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "astra_act_processing_duration_seconds",
            Help: "Time spent processing acts",
            Buckets: prometheus.DefBuckets,
        },
        []string{"act_type"},
    )
    
    // Business integration metrics
    integrationCallsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "astra_integration_calls_total", 
            Help: "Total number of business system integration calls",
        },
        []string{"system", "operation", "status"},
    )
    
    integrationDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "astra_integration_duration_seconds",
            Help: "Time spent on business system integrations",
            Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
        },
        []string{"system", "operation"},
    )
    
    // Conversation state metrics
    activeConversations = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "astra_active_conversations",
            Help: "Number of currently active conversations",
        },
    )
    
    // Error metrics
    errorsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "astra_errors_total",
            Help: "Total number of errors by type and system",
        },
        []string{"error_type", "system", "recoverable"},
    )
)

type MetricsCollector struct {
    // Implementation would include methods to record metrics
}

func (m *MetricsCollector) RecordActProcessed(actType, status string, duration float64) {
    actsProcessedTotal.WithLabelValues(actType, status).Inc()
    actProcessingDuration.WithLabelValues(actType).Observe(duration)
}

func (m *MetricsCollector) RecordIntegrationCall(system, operation, status string, duration float64) {
    integrationCallsTotal.WithLabelValues(system, operation, status).Inc()
    integrationDuration.WithLabelValues(system, operation).Observe(duration)
}

func (m *MetricsCollector) RecordError(errorType, system string, recoverable bool) {
    recoverableStr := "false"
    if recoverable {
        recoverableStr = "true"
    }
    errorsTotal.WithLabelValues(errorType, system, recoverableStr).Inc()
}
```

### Health Checks and Service Discovery

```typescript
class HealthCheckEndpoints {
  constructor(
    private conversationProcessor: ConversationProcessor,
    private integrationServices: Map<string, IntegrationService>
  ) {}

  // Basic liveness check
  async liveness(): Promise<HealthStatus> {
    return {
      status: 'healthy',
      timestamp: new Date().toISOString(),
      version: process.env.VERSION || 'unknown'
    };
  }

  // Comprehensive readiness check
  async readiness(): Promise<HealthStatus> {
    const checks: HealthCheck[] = [];

    // Check conversation processor
    checks.push(await this.checkConversationProcessor());
    
    // Check all integration services
    for (const [name, service] of this.integrationServices) {
      checks.push(await this.checkIntegrationService(name, service));
    }

    // Check external dependencies
    checks.push(await this.checkDatabase());
    checks.push(await this.checkMessageQueue());
    checks.push(await this.checkRedis());

    const allHealthy = checks.every(check => check.status === 'healthy');
    
    return {
      status: allHealthy ? 'healthy' : 'unhealthy',
      timestamp: new Date().toISOString(),
      checks
    };
  }

  private async checkConversationProcessor(): Promise<HealthCheck> {
    try {
      await this.conversationProcessor.healthCheck();
      return {
        name: 'conversation-processor',
        status: 'healthy',
        responseTime: Date.now()
      };
    } catch (error) {
      return {
        name: 'conversation-processor',
        status: 'unhealthy',
        error: error.message,
        responseTime: Date.now()
      };
    }
  }

  private async checkIntegrationService(name: string, service: IntegrationService): Promise<HealthCheck> {
    const startTime = Date.now();
    try {
      await service.healthCheck();
      return {
        name: `integration-${name}`,
        status: 'healthy',
        responseTime: Date.now() - startTime
      };
    } catch (error) {
      return {
        name: `integration-${name}`,
        status: 'unhealthy',
        error: error.message,
        responseTime: Date.now() - startTime
      };
    }
  }
}

interface HealthStatus {
  status: 'healthy' | 'unhealthy';
  timestamp: string;
  version?: string;
  checks?: HealthCheck[];
}

interface HealthCheck {
  name: string;
  status: 'healthy' | 'unhealthy';
  responseTime: number;
  error?: string;
}
```

## Security and Compliance

### Authentication and Authorization

```python
from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
import jwt
from typing import List, Optional

app = FastAPI()
security = HTTPBearer()

class SecurityManager:
    def __init__(self, secret_key: str, algorithm: str = "HS256"):
        self.secret_key = secret_key
        self.algorithm = algorithm
        
    async def verify_token(self, credentials: HTTPAuthorizationCredentials = Depends(security)):
        """Verify JWT token and extract user information"""
        try:
            payload = jwt.decode(
                credentials.credentials, 
                self.secret_key, 
                algorithms=[self.algorithm]
            )
            
            user_id = payload.get("sub")
            permissions = payload.get("permissions", [])
            organization = payload.get("organization")
            
            if not user_id:
                raise HTTPException(
                    status_code=status.HTTP_401_UNAUTHORIZED,
                    detail="Invalid token"
                )
                
            return {
                "user_id": user_id,
                "permissions": permissions,
                "organization": organization
            }
            
        except jwt.InvalidTokenError:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid token"
            )
    
    def require_permissions(self, required_permissions: List[str]):
        """Decorator to require specific permissions"""
        def permission_checker(user = Depends(self.verify_token)):
            user_permissions = set(user.get("permissions", []))
            required_permissions_set = set(required_permissions)
            
            if not required_permissions_set.issubset(user_permissions):
                raise HTTPException(
                    status_code=status.HTTP_403_FORBIDDEN,
                    detail="Insufficient permissions"
                )
            return user
        return permission_checker

security_manager = SecurityManager(secret_key="your-secret-key")

@app.post("/conversation/{conversation_id}/acts")
async def submit_act(
    conversation_id: str,
    act: dict,
    user = Depends(security_manager.require_permissions(["conversation.write"]))
):
    """Submit a new act to a conversation (requires write permission)"""
    
    # Verify user has access to this conversation
    if not await verify_conversation_access(conversation_id, user["user_id"], user["organization"]):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Access denied to conversation"
        )
    
    # Process the act
    return await process_conversation_act(conversation_id, act, user)

@app.get("/conversation/{conversation_id}")
async def get_conversation(
    conversation_id: str,
    user = Depends(security_manager.require_permissions(["conversation.read"]))
):
    """Get conversation details (requires read permission)"""
    
    if not await verify_conversation_access(conversation_id, user["user_id"], user["organization"]):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Access denied to conversation"
        )
    
    return await fetch_conversation(conversation_id, user)

async def verify_conversation_access(
    conversation_id: str, 
    user_id: str, 
    organization: str
) -> bool:
    """Verify user has access to conversation based on organization and permissions"""
    # Implementation would check conversation ownership/access rules
    pass
```

### Data Privacy and Compliance

```typescript
class PrivacyManager {
  private piiDetector: PIIDetector;
  private encryptionService: EncryptionService;

  constructor(
    piiDetector: PIIDetector,
    encryptionService: EncryptionService
  ) {
    this.piiDetector = piiDetector;
    this.encryptionService = encryptionService;
  }

  async processActForPrivacy(act: ConversationAct, privacyPolicy: PrivacyPolicy): Promise<ConversationAct> {
    // Detect PII in the act
    const piiElements = await this.piiDetector.detect(act);
    
    if (piiElements.length === 0) {
      return act; // No PII detected
    }

    // Apply privacy policy
    return await this.applyPrivacyPolicy(act, piiElements, privacyPolicy);
  }

  private async applyPrivacyPolicy(
    act: ConversationAct,
    piiElements: PIIElement[],
    policy: PrivacyPolicy
  ): Promise<ConversationAct> {
    let processedAct = { ...act };

    for (const pii of piiElements) {
      const action = policy.getPolicyForPIIType(pii.type);
      
      switch (action) {
        case 'encrypt':
          processedAct = await this.encryptPIIInAct(processedAct, pii);
          break;
        case 'hash':
          processedAct = await this.hashPIIInAct(processedAct, pii);
          break;
        case 'redact':
          processedAct = await this.redactPIIInAct(processedAct, pii);
          break;
        case 'tokenize':
          processedAct = await this.tokenizePIIInAct(processedAct, pii);
          break;
      }
    }

    // Add privacy metadata
    processedAct.metadata = {
      ...processedAct.metadata,
      privacy_processed: true,
      pii_elements_count: piiElements.length,
      processing_timestamp: new Date().toISOString()
    };

    return processedAct;
  }

  private async encryptPIIInAct(act: ConversationAct, pii: PIIElement): Promise<ConversationAct> {
    const encryptedValue = await this.encryptionService.encrypt(pii.value);
    return this.replacePIIInAct(act, pii, `[ENCRYPTED:${encryptedValue.id}]`);
  }

  private async redactPIIInAct(act: ConversationAct, pii: PIIElement): Promise<ConversationAct> {
    const redactedValue = '[REDACTED:' + pii.type.toUpperCase() + ']';
    return this.replacePIIInAct(act, pii, redactedValue);
  }

  private replacePIIInAct(act: ConversationAct, pii: PIIElement, replacement: string): ConversationAct {
    // Deep clone and replace PII based on field path
    const processedAct = JSON.parse(JSON.stringify(act));
    this.setNestedProperty(processedAct, pii.fieldPath, replacement);
    return processedAct;
  }
}

interface PIIElement {
  type: 'email' | 'phone' | 'ssn' | 'credit_card' | 'address' | 'name';
  value: string;
  fieldPath: string[];
  confidence: number;
}

interface PrivacyPolicy {
  getPolicyForPIIType(type: string): 'encrypt' | 'hash' | 'redact' | 'tokenize' | 'allow';
}
```

## Testing Integration Implementations

### Integration Testing Framework

```python
import pytest
from unittest.mock import AsyncMock, patch
from astra_model import Ask, Fact, Commit, Conversation

class TestCRMIntegration:
    @pytest.fixture
    def crm_integration(self):
        mock_crm_client = AsyncMock()
        return CRMIntegration(mock_crm_client)
    
    @pytest.fixture
    def sample_customer_fact(self):
        return Fact(
            id="act_001",
            timestamp="2025-01-15T14:30:00Z",
            speaker="customer_123",
            type="fact",
            entity={
                "id": "customer_456",
                "type": "customer",
                "external_id": "crm_789"
            },
            field="email",
            value="updated@example.com"
        )
    
    @pytest.mark.asyncio
    async def test_customer_fact_updates_crm(self, crm_integration, sample_customer_fact):
        """Test that customer facts update CRM records correctly"""
        
        # Act
        await crm_integration.handleCustomerFacts([sample_customer_fact])
        
        # Assert
        crm_integration.crmClient.updateContact.assert_called_once_with(
            "crm_789",  # external_id used as CRM ID
            {"email_address": "updated@example.com"}
        )
    
    @pytest.mark.asyncio
    async def test_crm_api_failure_handling(self, crm_integration, sample_customer_fact):
        """Test error handling when CRM API fails"""
        
        # Arrange
        crm_integration.crmClient.updateContact.side_effect = Exception("CRM API Error")
        
        # Act & Assert
        with pytest.raises(Exception, match="CRM API Error"):
            await crm_integration.handleCustomerFacts([sample_customer_fact])
        
        # Verify error was recorded
        assert crm_integration.error_count > 0
    
    @pytest.mark.asyncio
    async def test_commit_creates_customer_record(self, crm_integration):
        """Test that customer commit creates new CRM record"""
        
        commit = Commit(
            id="act_002",
            timestamp="2025-01-15T14:31:00Z",
            speaker="system",
            type="commit",
            entity="customer_456",
            action="create",
            system="crm"
        )
        
        # Mock successful creation
        crm_integration.crmClient.createContact.return_value = {"id": "crm_new_123"}
        
        # Act
        await crm_integration.handleCustomerCommits(commit)
        
        # Assert
        crm_integration.crmClient.createContact.assert_called_once()

class TestOrderIntegrationWorkflow:
    @pytest.fixture
    def order_integration(self):
        mock_order_service = AsyncMock()
        mock_inventory_service = AsyncMock()
        mock_payment_service = AsyncMock()
        
        return OrderManagementIntegration(
            mock_order_service,
            mock_inventory_service,
            mock_payment_service
        )
    
    @pytest.fixture
    def complete_order_state(self):
        return {
            "customer_id": "customer_123",
            "items": [
                {"product_id": "pizza_large", "quantity": 2, "price": 15.99}
            ],
            "total_amount": 31.98,
            "delivery_address": "123 Main St, Anytown, USA",
            "payment_method": {
                "type": "credit_card",
                "card_number": "4111111111111111"
            }
        }
    
    @pytest.mark.asyncio
    async def test_successful_order_creation_workflow(self, order_integration, complete_order_state):
        """Test complete order creation workflow with all validations"""
        
        # Arrange
        commit = Commit(
            id="act_003",
            timestamp="2025-01-15T14:32:00Z",
            speaker="system",
            type="commit",
            entity="order_789",
            action="create",
            system="order_management"
        )
        
        # Mock successful responses
        order_integration.inventory_service.check_availability.return_value = True
        order_integration.inventory_service.reserve_inventory.return_value = "reservation_123"
        order_integration.order_service.create_order.return_value = {"id": "order_789"}
        order_integration.payment_service.charge.return_value = {
            "status": "success",
            "transaction_id": "payment_456"
        }
        
        # Act
        await order_integration.handle_order_commit(commit, complete_order_state)
        
        # Assert
        order_integration.inventory_service.check_availability.assert_called()
        order_integration.inventory_service.reserve_inventory.assert_called()
        order_integration.order_service.create_order.assert_called()
        order_integration.payment_service.charge.assert_called()
    
    @pytest.mark.asyncio
    async def test_inventory_shortage_rollback(self, order_integration, complete_order_state):
        """Test that inventory reservations are rolled back on failure"""
        
        commit = Commit(
            id="act_004",
            timestamp="2025-01-15T14:33:00Z",
            speaker="system",
            type="commit",
            entity="order_789",
            action="create",
            system="order_management"
        )
        
        # Arrange - inventory available but payment fails
        order_integration.inventory_service.check_availability.return_value = True
        order_integration.inventory_service.reserve_inventory.return_value = "reservation_123"
        order_integration.payment_service.charge.side_effect = Exception("Payment failed")
        
        # Act & Assert
        with pytest.raises(Exception, match="Payment failed"):
            await order_integration.handle_order_commit(commit, complete_order_state)
        
        # Verify rollback
        order_integration.inventory_service.cancel_reservation.assert_called_with("reservation_123")

# Load testing for high-throughput scenarios
class TestHighThroughputIntegration:
    @pytest.mark.asyncio
    async def test_concurrent_act_processing(self):
        """Test processing multiple acts concurrently"""
        
        processor = ConversationStreamProcessor(test_kafka_config)
        
        # Create 100 concurrent acts
        acts = [
            create_test_fact(f"act_{i}", f"entity_{i}", "field", f"value_{i}")
            for i in range(100)
        ]
        
        # Process concurrently
        start_time = time.time()
        
        tasks = [processor.process_message({"act": act}) for act in acts]
        await asyncio.gather(*tasks)
        
        end_time = time.time()
        processing_time = end_time - start_time
        
        # Assert reasonable performance (< 5 seconds for 100 acts)
        assert processing_time < 5.0
        assert processor.processed_count == 100

def create_test_fact(act_id: str, entity_id: str, field: str, value: str) -> dict:
    """Helper function to create test facts"""
    return {
        "id": act_id,
        "timestamp": datetime.utcnow().isoformat() + "Z",
        "speaker": "test_speaker",
        "type": "fact",
        "entity": entity_id,
        "field": field,
        "value": value
    }
```

## Performance Optimization

### Batch Processing Strategies

```typescript
class BatchProcessor {
  private batchSize: number = 100;
  private batchTimeout: number = 5000; // 5 seconds
  private pendingActs: ConversationAct[] = [];
  private batchTimer?: NodeJS.Timeout;

  constructor(
    private integrationService: IntegrationService,
    batchSize = 100,
    batchTimeout = 5000
  ) {
    this.batchSize = batchSize;
    this.batchTimeout = batchTimeout;
  }

  async addAct(act: ConversationAct): Promise<void> {
    this.pendingActs.push(act);

    // Process batch if it's full
    if (this.pendingActs.length >= this.batchSize) {
      await this.processBatch();
      return;
    }

    // Set timer for timeout-based processing
    if (!this.batchTimer) {
      this.batchTimer = setTimeout(
        () => this.processBatch(),
        this.batchTimeout
      );
    }
  }

  private async processBatch(): Promise<void> {
    if (this.pendingActs.length === 0) return;

    const batch = [...this.pendingActs];
    this.pendingActs = [];

    // Clear timeout
    if (this.batchTimer) {
      clearTimeout(this.batchTimer);
      this.batchTimer = undefined;
    }

    try {
      await this.processBatchWithRetry(batch);
    } catch (error) {
      // Handle batch processing error
      await this.handleBatchError(batch, error);
    }
  }

  private async processBatchWithRetry(
    batch: ConversationAct[],
    maxRetries = 3
  ): Promise<void> {
    let lastError;

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        await this.integrationService.processBatch(batch);
        return; // Success
      } catch (error) {
        lastError = error;
        
        if (attempt < maxRetries) {
          const delay = Math.min(1000 * Math.pow(2, attempt), 10000);
          await new Promise(resolve => setTimeout(resolve, delay));
        }
      }
    }

    throw new BatchProcessingError(
      `Failed to process batch after ${maxRetries + 1} attempts`,
      lastError,
      batch
    );
  }

  // Graceful shutdown
  async shutdown(): Promise<void> {
    if (this.batchTimer) {
      clearTimeout(this.batchTimer);
    }
    
    if (this.pendingActs.length > 0) {
      await this.processBatch();
    }
  }
}
```

## Best Practices Summary

### Integration Architecture
1. **Use event-driven patterns** for scalability and loose coupling
2. **Implement circuit breakers** for external service calls
3. **Design for idempotency** to handle retries safely
4. **Separate concerns** between conversation processing and business integration
5. **Use batch processing** for high-throughput scenarios

### Error Handling
1. **Classify errors** as retryable vs. non-retryable
2. **Implement exponential backoff** with jitter for retries
3. **Use dead letter queues** for failed messages
4. **Monitor error rates** and alert on anomalies
5. **Provide graceful degradation** when services are unavailable

### Security and Privacy
1. **Encrypt sensitive data** at rest and in transit
2. **Implement proper authentication** and authorization
3. **Audit all data access** and modifications
4. **Handle PII appropriately** based on privacy policies
5. **Use secure communication channels** for all integrations

### Monitoring and Operations
1. **Implement comprehensive logging** with structured data
2. **Use distributed tracing** for end-to-end visibility
3. **Monitor business metrics** alongside technical metrics
4. **Set up proper alerting** for critical failures
5. **Provide health check endpoints** for all services

### Performance and Scalability
1. **Design for horizontal scaling** from the start
2. **Use connection pooling** for database and external services
3. **Implement caching** where appropriate
4. **Monitor resource usage** and optimize bottlenecks
5. **Load test integrations** under realistic conditions

## Migration Strategies

### From Legacy Systems

When migrating from existing conversational systems to ASTRA, a phased approach minimizes risk and enables gradual adoption.

```typescript
// Legacy system adapter pattern
class LegacySystemAdapter {
  constructor(
    private legacyAPI: LegacyConversationAPI,
    private astraProcessor: ASTRAProcessor
  ) {}

  async processLegacyEvent(legacyEvent: LegacyEvent): Promise<void> {
    try {
      // Transform legacy event to ASTRA acts
      const acts = await this.transformToASTRA(legacyEvent);
      
      // Process each act through ASTRA pipeline
      for (const act of acts) {
        await this.astraProcessor.processAct(act);
      }
      
      // Maintain backward compatibility by updating legacy system
      await this.updateLegacySystem(legacyEvent, acts);
      
    } catch (error) {
      // Fallback to legacy processing on failure
      await this.legacyFallback(legacyEvent, error);
    }
  }

  private async transformToASTRA(legacyEvent: LegacyEvent): Promise<ConversationAct[]> {
    const acts: ConversationAct[] = [];
    
    // Map legacy intents to ASTRA acts
    switch (legacyEvent.intent) {
      case 'collect_information':
        acts.push(this.createAskAct(legacyEvent));
        break;
        
      case 'provide_information':
        acts.push(this.createFactAct(legacyEvent));
        break;
        
      case 'confirm_details':
        acts.push(this.createConfirmAct(legacyEvent));
        break;
        
      case 'execute_action':
        acts.push(this.createCommitAct(legacyEvent));
        break;
        
      default:
        throw new UnknownIntentError(`Unknown intent: ${legacyEvent.intent}`);
    }
    
    return acts;
  }

  private createFactAct(legacyEvent: LegacyEvent): Fact {
    return {
      id: `act_${legacyEvent.id}`,
      timestamp: legacyEvent.timestamp,
      speaker: legacyEvent.userId,
      type: 'fact',
      entity: legacyEvent.entityId || 'unknown',
      field: legacyEvent.field,
      value: legacyEvent.value,
      metadata: {
        legacy_source: true,
        original_intent: legacyEvent.intent,
        confidence: legacyEvent.confidence
      }
    };
  }

  private async legacyFallback(legacyEvent: LegacyEvent, error: Error): Promise<void> {
    // Log the transformation failure
    console.error('ASTRA transformation failed, falling back to legacy processing', {
      eventId: legacyEvent.id,
      error: error.message
    });
    
    // Process using legacy system
    await this.legacyAPI.processEvent(legacyEvent);
    
    // Record metrics for monitoring migration progress
    this.recordMigrationFailure(legacyEvent, error);
  }
}
```

### Gradual Feature Migration

```python
class FeatureFlaggedProcessor:
    """Enables gradual migration using feature flags"""
    
    def __init__(self, feature_flags, legacy_processor, astra_processor):
        self.feature_flags = feature_flags
        self.legacy_processor = legacy_processor
        self.astra_processor = astra_processor
    
    async def process_conversation(self, conversation_data):
        """Route processing based on feature flags"""
        
        # Check if ASTRA is enabled for this conversation/user/organization
        use_astra = await self.should_use_astra(conversation_data)
        
        if use_astra:
            try:
                return await self.astra_processor.process(conversation_data)
            except Exception as e:
                # Fallback to legacy on ASTRA failure
                if self.feature_flags.get('astra_failover_enabled'):
                    self.log_astra_failure(e, conversation_data)
                    return await self.legacy_processor.process(conversation_data)
                raise e
        else:
            return await self.legacy_processor.process(conversation_data)
    
    async def should_use_astra(self, conversation_data) -> bool:
        """Determine if ASTRA should be used for this conversation"""
        
        # Check global rollout percentage
        rollout_percentage = self.feature_flags.get('astra_rollout_percentage', 0)
        if random.randint(1, 100) > rollout_percentage:
            return False
        
        # Check user-specific flags
        user_id = conversation_data.get('user_id')
        if user_id in self.feature_flags.get('astra_beta_users', []):
            return True
        
        # Check organization-specific flags
        org_id = conversation_data.get('organization_id')
        if org_id in self.feature_flags.get('astra_enabled_orgs', []):
            return True
        
        # Check conversation type
        conversation_type = conversation_data.get('type')
        enabled_types = self.feature_flags.get('astra_enabled_conversation_types', [])
        
        return conversation_type in enabled_types
```

## Troubleshooting and Debugging

### Common Integration Issues

```typescript
class IntegrationDiagnostics {
  
  async diagnoseIntegrationFailure(
    act: ConversationAct,
    error: Error,
    context: IntegrationContext
  ): Promise<DiagnosticReport> {
    const report: DiagnosticReport = {
      act_id: act.id,
      error_type: error.constructor.name,
      error_message: error.message,
      timestamp: new Date().toISOString(),
      checks: []
    };

    // Check act validation
    report.checks.push(await this.checkActValidation(act));
    
    // Check entity state
    if (isFactOrCommit(act)) {
      report.checks.push(await this.checkEntityState(act.entity));
    }
    
    // Check external system connectivity
    if (isCommit(act) && act.system) {
      report.checks.push(await this.checkSystemConnectivity(act.system));
    }
    
    // Check authentication
    report.checks.push(await this.checkAuthentication(context));
    
    // Check rate limits
    report.checks.push(await this.checkRateLimits(context));
    
    return report;
  }

  private async checkActValidation(act: ConversationAct): Promise<DiagnosticCheck> {
    try {
      const validator = new ActValidator();
      const result = await validator.validate(act);
      
      return {
        name: 'act_validation',
        status: result.valid ? 'pass' : 'fail',
        details: result.valid ? 'Act structure is valid' : result.errors,
        suggestions: result.valid ? [] : [
          'Check required fields are present',
          'Verify field types match schema',
          'Ensure enum values are valid'
        ]
      };
    } catch (error) {
      return {
        name: 'act_validation',
        status: 'error',
        details: `Validation check failed: ${error.message}`,
        suggestions: ['Check act validator configuration']
      };
    }
  }

  private async checkSystemConnectivity(system: string): Promise<DiagnosticCheck> {
    try {
      const integration = this.getIntegrationForSystem(system);
      await integration.healthCheck();
      
      return {
        name: 'system_connectivity',
        status: 'pass',
        details: `Connection to ${system} is healthy`
      };
    } catch (error) {
      return {
        name: 'system_connectivity',
        status: 'fail',
        details: `Cannot connect to ${system}: ${error.message}`,
        suggestions: [
          'Check system is running and accessible',
          'Verify network connectivity',
          'Check authentication credentials',
          'Review rate limiting settings'
        ]
      };
    }
  }
}

interface DiagnosticReport {
  act_id: string;
  error_type: string;
  error_message: string;
  timestamp: string;
  checks: DiagnosticCheck[];
}

interface DiagnosticCheck {
  name: string;
  status: 'pass' | 'fail' | 'error';
  details: string;
  suggestions?: string[];
}
```

### Debug Logging and Tracing

```python
import logging
import json
from typing import Dict, Any
from contextvars import ContextVar

# Context variable for tracing
trace_context: ContextVar[Dict[str, Any]] = ContextVar('trace_context', default={})

class ASTRALogger:
    """Enhanced logging for ASTRA integrations"""
    
    def __init__(self, name: str):
        self.logger = logging.getLogger(name)
        self.logger.setLevel(logging.INFO)
        
        # Create structured formatter
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
        
        handler = logging.StreamHandler()
        handler.setFormatter(formatter)
        self.logger.addHandler(handler)
    
    def log_act_processing(self, act: dict, stage: str, **kwargs):
        """Log act processing with full context"""
        context = trace_context.get({})
        
        log_data = {
            'stage': stage,
            'act_id': act.get('id'),
            'act_type': act.get('type'),
            'conversation_id': context.get('conversation_id'),
            'trace_id': context.get('trace_id'),
            **kwargs
        }
        
        self.logger.info(f"Act processing: {stage}", extra=log_data)
    
    def log_integration_call(
        self, 
        system: str, 
        operation: str, 
        duration_ms: float,
        success: bool,
        **kwargs
    ):
        """Log external system integration calls"""
        context = trace_context.get({})
        
        log_data = {
            'integration_system': system,
            'integration_operation': operation,
            'duration_ms': duration_ms,
            'success': success,
            'conversation_id': context.get('conversation_id'),
            'trace_id': context.get('trace_id'),
            **kwargs
        }
        
        level = logging.INFO if success else logging.ERROR
        self.logger.log(level, f"Integration call: {system}.{operation}", extra=log_data)
    
    def log_error(self, error: Exception, context_info: dict = None):
        """Log errors with full context"""
        context = trace_context.get({})
        
        log_data = {
            'error_type': type(error).__name__,
            'error_message': str(error),
            'conversation_id': context.get('conversation_id'),
            'trace_id': context.get('trace_id'),
            **(context_info or {})
        }
        
        self.logger.error("Integration error occurred", extra=log_data, exc_info=True)

# Usage example
class TracedOrderIntegration:
    def __init__(self):
        self.logger = ASTRALogger('order_integration')
    
    async def process_order_commit(self, commit: dict, conversation_id: str):
        # Set trace context
        trace_context.set({
            'conversation_id': conversation_id,
            'trace_id': generate_trace_id(),
            'integration': 'order_management'
        })
        
        self.logger.log_act_processing(commit, 'validation_start')
        
        try:
            # Validate commit
            await self.validate_commit(commit)
            self.logger.log_act_processing(commit, 'validation_complete')
            
            # Execute integration
            start_time = time.time()
            result = await self.execute_order_action(commit)
            duration = (time.time() - start_time) * 1000
            
            self.logger.log_integration_call(
                'order_management',
                commit['action'],
                duration,
                True,
                transaction_id=result.get('transaction_id')
            )
            
        except Exception as e:
            self.logger.log_error(e, {
                'commit_id': commit['id'],
                'entity_id': commit.get('entity')
            })
            raise
```

## Advanced Patterns

### Saga Pattern for Distributed Transactions

```typescript
// Implementing the Saga pattern for complex multi-system commits
class OrderSagaOrchestrator {
  private sagaSteps: SagaStep[] = [
    { name: 'validate_order', compensate: 'noop' },
    { name: 'reserve_inventory', compensate: 'release_inventory' },
    { name: 'process_payment', compensate: 'refund_payment' },
    { name: 'create_order', compensate: 'cancel_order' },
    { name: 'update_crm', compensate: 'revert_crm_update' },
    { name: 'send_confirmation', compensate: 'send_cancellation' }
  ];

  async executeOrderSaga(
    orderCommit: Commit,
    orderState: OrderEntity
  ): Promise<SagaResult> {
    const sagaId = `saga_${orderCommit.id}_${Date.now()}`;
    const executedSteps: string[] = [];

    try {
      for (const step of this.sagaSteps) {
        await this.executeStep(step.name, orderCommit, orderState);
        executedSteps.push(step.name);
        
        // Record saga progress
        await this.recordSagaProgress(sagaId, step.name, 'completed');
      }

      await this.recordSagaCompletion(sagaId, 'success');
      return { status: 'success', sagaId };

    } catch (error) {
      // Execute compensating actions in reverse order
      await this.compensate(executedSteps.reverse(), orderCommit, orderState);
      await this.recordSagaCompletion(sagaId, 'failed', error);
      
      throw new SagaFailureError(
        `Saga ${sagaId} failed at step ${executedSteps[0]}`,
        error
      );
    }
  }

  private async compensate(
    executedSteps: string[],
    orderCommit: Commit,
    orderState: OrderEntity
  ): Promise<void> {
    for (const stepName of executedSteps) {
      try {
        const step = this.sagaSteps.find(s => s.name === stepName);
        if (step && step.compensate !== 'noop') {
          await this.executeCompensation(step.compensate, orderCommit, orderState);
        }
      } catch (compensationError) {
        // Log compensation failures but continue with other compensations
        console.error(`Compensation failed for step ${stepName}`, compensationError);
      }
    }
  }

  private async executeStep(
    stepName: string,
    commit: Commit,
    state: OrderEntity
  ): Promise<void> {
    const stepHandlers = {
      validate_order: () => this.validateOrder(state),
      reserve_inventory: () => this.reserveInventory(state),
      process_payment: () => this.processPayment(state),
      create_order: () => this.createOrder(state),
      update_crm: () => this.updateCRM(state),
      send_confirmation: () => this.sendConfirmation(state)
    };

    const handler = stepHandlers[stepName];
    if (!handler) {
      throw new Error(`Unknown saga step: ${stepName}`);
    }

    await handler();
  }
}

interface SagaStep {
  name: string;
  compensate: string;
}

interface SagaResult {
  status: 'success' | 'failed';
  sagaId: string;
  error?: Error;
}
```

### Event Sourcing Integration

```python
from typing import List, Dict, Any
from datetime import datetime
import json

class EventStore:
    """Event store for ASTRA conversation events"""
    
    def __init__(self, storage_backend):
        self.storage = storage_backend
    
    async def append_events(
        self, 
        stream_id: str, 
        events: List[Dict[str, Any]],
        expected_version: int = -1
    ) -> int:
        """Append events to a stream with optimistic concurrency control"""
        
        # Check expected version for concurrency control
        current_version = await self.get_stream_version(stream_id)
        if expected_version != -1 and current_version != expected_version:
            raise ConcurrencyError(
                f"Expected version {expected_version}, got {current_version}"
            )
        
        # Prepare events with metadata
        enriched_events = []
        for i, event in enumerate(events):
            enriched_events.append({
                'stream_id': stream_id,
                'event_number': current_version + i + 1,
                'event_type': event['type'],
                'data': event,
                'metadata': {
                    'timestamp': datetime.utcnow().isoformat(),
                    'correlation_id': event.get('correlation_id'),
                    'causation_id': event.get('causation_id')
                }
            })
        
        # Atomic append
        await self.storage.append_events(stream_id, enriched_events)
        return current_version + len(events)

class ConversationProjection:
    """Projection that builds conversation state from events"""
    
    def __init__(self, event_store: EventStore):
        self.event_store = event_store
        self.state_cache = {}
    
    async def get_conversation_state(self, conversation_id: str) -> Dict[str, Any]:
        """Get current conversation state by replaying events"""
        
        if conversation_id in self.state_cache:
            return self.state_cache[conversation_id]
        
        # Load events from event store
        events = await self.event_store.get_stream_events(
            f"conversation-{conversation_id}"
        )
        
        # Replay events to build state
        state = self.replay_events(events)
        
        # Cache the state
        self.state_cache[conversation_id] = state
        return state
    
    def replay_events(self, events: List[Dict[str, Any]]) -> Dict[str, Any]:
        """Replay events to build current state"""
        
        state = {
            'participants': {},
            'entities': {},
            'acts': [],
            'status': 'active'
        }
        
        for event in events:
            self.apply_event(state, event)
        
        return state
    
    def apply_event(self, state: Dict[str, Any], event: Dict[str, Any]):
        """Apply a single event to the state"""
        
        event_data = event['data']
        event_type = event['event_type']
        
        if event_type == 'ActProcessed':
            act = event_data['act']
            state['acts'].append(act)
            
            # Update entity state based on facts
            if act['type'] == 'fact':
                entity_id = act['entity']
                if entity_id not in state['entities']:
                    state['entities'][entity_id] = {}
                
                state['entities'][entity_id][act['field']] = act['value']
        
        elif event_type == 'ParticipantJoined':
            participant = event_data['participant']
            state['participants'][participant['id']] = participant
        
        elif event_type == 'ConversationCompleted':
            state['status'] = 'completed'
        
        elif event_type == 'ConversationFailed':
            state['status'] = 'failed'
            state['error'] = event_data['error']

class EventSourcingIntegration:
    """Integration that uses event sourcing for ASTRA conversations"""
    
    def __init__(self, event_store: EventStore, projection: ConversationProjection):
        self.event_store = event_store
        self.projection = projection
    
    async def process_act(self, act: dict, conversation_id: str):
        """Process act using event sourcing pattern"""
        
        try:
            # Get current state
            current_state = await self.projection.get_conversation_state(conversation_id)
            
            # Process the act
            result = await self.apply_act_to_state(act, current_state)
            
            # Generate events based on processing result
            events = self.generate_events(act, result)
            
            # Append events to stream
            stream_id = f"conversation-{conversation_id}"
            await self.event_store.append_events(stream_id, events)
            
            # Invalidate cache to force reload
            if conversation_id in self.projection.state_cache:
                del self.projection.state_cache[conversation_id]
            
            return result
            
        except Exception as e:
            # Generate error event
            error_event = {
                'type': 'ActProcessingFailed',
                'act_id': act['id'],
                'error': str(e),
                'correlation_id': act.get('id')
            }
            
            await self.event_store.append_events(
                f"conversation-{conversation_id}",
                [error_event]
            )
            raise
    
    def generate_events(self, act: dict, result: dict) -> List[Dict[str, Any]]:
        """Generate events based on act processing result"""
        events = []
        
        # Always generate act processed event
        events.append({
            'type': 'ActProcessed',
            'act': act,
            'result': result,
            'correlation_id': act['id']
        })
        
        # Generate additional events based on act type
        if act['type'] == 'commit':
            if result.get('success'):
                events.append({
                    'type': 'BusinessActionExecuted',
                    'commit_id': act['id'],
                    'action': act['action'],
                    'entity': act['entity'],
                    'transaction_id': result.get('transaction_id'),
                    'causation_id': act['id']
                })
            else:
                events.append({
                    'type': 'BusinessActionFailed',
                    'commit_id': act['id'],
                    'error': result.get('error'),
                    'causation_id': act['id']
                })
        
        return events
```

This comprehensive integration guide provides everything needed to successfully implement ASTRA in production environments. The patterns and practices outlined here enable robust, scalable, and maintainable conversational applications that can evolve with changing business requirements while maintaining high reliability and performance standards.
