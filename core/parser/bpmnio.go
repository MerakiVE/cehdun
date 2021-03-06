package parser

import (
	"strings"
	"github.com/beevik/etree"
	"github.com/merakiVE/koinos/utils"
)

/*
  'bpmn:Association',
  'bpmn:BusinessRuleTask',
  'bpmn:DataInputAssociation',
  'bpmn:DataOutputAssociation',
  'bpmn:DataObjectReference',
  'bpmn:DataStoreReference',
  'bpmn:EndEvent',
  'bpmn:EventBasedGateway',
  'bpmn:ExclusiveGateway',
  'bpmn:IntermediateCatchEvent',
  'bpmn:ManualTask',
  'bpmn:ParallelGateway',
  'bpmn:Process',
  'bpmn:SequenceFlow',
  'bpmn:StartEvent',
  'bpmn:SubProcess',
  'bpmn:Task',
  'bpmn:TextAnnotation',
  'bpmn:UserTask'
*/

const (
	EMPTY     = "unknown"
	XML_SPACE = "bpmn2"

	BPMNIO_TYPE_GATEWAY  = "gateway"
	BPMNIO_TYPE_EVENT    = "event"
	BPMNIO_TYPE_TASK     = "task"
	BPMNIO_TYPE_ACTIVITY = "activity"
	BPMNIO_TYPE_NONE     = ""

	BPMNIO_TAG_ROOT                    = XML_SPACE + ":definitions"
	BPMNIO_TAG_COLABORATION            = XML_SPACE + ":collaboration"
	BPMNIO_TAG_PROCESS                 = XML_SPACE + ":process"
	BPMNIO_TAG_SUB_PROCESS             = XML_SPACE + ":subProcess"
	BPMNIO_TAG_START_EVENT             = XML_SPACE + ":startEvent"
	BPMNIO_TAG_END_EVENT               = XML_SPACE + ":endEvent"
	BPMNIO_TAG_OUTGOING                = XML_SPACE + ":outgoing"
	BPMNIO_TAG_INCOMING                = XML_SPACE + ":incoming"
	BPMNIO_TAG_SEQUENCE_FLOW           = XML_SPACE + ":sequenceFlow"
	BPMNIO_TAG_MESSAGE_FLOW            = XML_SPACE + ":messageFlow"
	BPMNIO_TAG_TASK                    = XML_SPACE + ":task"
	BPMNIO_TAG_LANE_SET                = XML_SPACE + ":laneSet"
	BPMNIO_TAG_LANE                    = XML_SPACE + ":lane"
	BPMNIO_TAG_FLOW_NODE_REF           = XML_SPACE + ":flowNodeRef"
	BPMNIO_TAG_EVENT_BASED_GATEWAY     = XML_SPACE + ":eventBasedGateway"
	BPMNIO_TAG_EVENT_EXCLUSIVE_GATEWAY = XML_SPACE + ":exclusiveGateway"
	BPMNIO_TAG_EVENT_PARALLEL_GATEWAY  = XML_SPACE + ":parallelGateway"
	BPMNIO_TAG_EVENT_COMPLEX_GATWAY    = XML_SPACE + ":complexGateway"
	BPMNIO_TAG_DATA_INPUT_ASSOCIATION  = XML_SPACE + ":dataInputAssociation"
	BPMNIO_TAG_DATA_OUTPUT_ASSOCIATION = XML_SPACE + ":dataOutputAssociation"
	BPMNIO_TAG_DATA_OBJECT_REFERENCE   = XML_SPACE + ":dataObjectReference"

	BPMNIO_ATTR_ID         = "id"
	BPMNIO_ATTR_SOURCE_REF = "sourceRef"
	BPMNIO_ATTR_TARGET_REF = "targetRef"
	BPMNIO_ATTR_NAME       = "name"

	BPMNIO_ATTR_CVDI_NEURON = "cvdi:neuron"
	BPMNIO_ATTR_CVDI_ACTION = "cvdi:action"
)

/* Funcion que crea la estructura de datos para el diagrama bpmn */

func NewParserBPMNIO() *DiagramBpmnIO {
	doc := etree.NewDocument()
	return &DiagramBpmnIO{documentXML: doc, flows: make([]*etree.Element, 0)}
}

/*
	Estructura para el diagrama BPMN de la herramienta bpmn.io
 */

type DiagramBpmnIO struct {
	documentXML *etree.Document
	flows       []*etree.Element
}

/* Funcion que carga el diagrama XML en forma de pathfile */
func (this *DiagramBpmnIO) ReadFromFile(filename string) error {
	this.documentXML = etree.NewDocument()
	if err := this.documentXML.ReadFromFile(filename); err != nil {
		return err
	}
	this.findAndLoadFlows()
	return nil
}

/* Funcion que carga el diagrama XML en forma de string */
func (this *DiagramBpmnIO) ReadFromString(data string) error {
	this.documentXML = etree.NewDocument()
	if err := this.documentXML.ReadFromString(data); err != nil {
		return err
	}
	this.findAndLoadFlows()
	return nil
}

/* Funcion que carga el diagrama XML en forma de byte */
func (this *DiagramBpmnIO) ReadFromBytes(bytes []byte) error {
	this.documentXML = etree.NewDocument()
	if err := this.documentXML.ReadFromBytes(bytes); err != nil {
		return err
	}
	this.findAndLoadFlows()
	return nil
}

/* Funcion que carga todos los flows en un slice de la estructura */
func (this *DiagramBpmnIO) findAndLoadFlows() {
	//Buscar todos los elementos padres que tengan una etiqueta TAG_MESSAGE_FLOW y TAG_SEQUENCE_FLOW como hijos
	parent_messages := this.getRootElement().FindElements(`[` + BPMNIO_TAG_MESSAGE_FLOW + `]`)
	parent_sequences := this.getRootElement().FindElements(`[` + BPMNIO_TAG_SEQUENCE_FLOW + `]`)

	//Anadimos todos los TAG_SEQUENCE_FLOW en el diagrama al atributo this.flows
	for _, mesg := range parent_messages {
		this.flows = append(this.flows, mesg.SelectElements(BPMNIO_TAG_MESSAGE_FLOW)...)
	}

	//Anadimos todos los TAG_SEQUENCE_FLOW en el diagrama al atributo this.flows
	for _, seq := range parent_sequences {
		this.flows = append(this.flows, seq.SelectElements(BPMNIO_TAG_SEQUENCE_FLOW)...)
	}
}

/* Verifica si un elemento es una estructura gateway */
func (this DiagramBpmnIO) isGateway(elem *etree.Element) (bool) {
	return strings.HasSuffix(elem.Tag, "Gateway")
}

/* Verifica si un elemento es un evento */
func (this DiagramBpmnIO) isEvent(elem *etree.Element) (bool) {
	return strings.HasSuffix(elem.Tag, "Event")
}

/* Verifica si un elemento es una tarea */
func (this DiagramBpmnIO) isTask(elem *etree.Element) (bool) {
	taks := []string{"subProcess", "transaction", "task"}
	for _, v := range taks {
		if v == elem.Tag {
			return true
		}
	}
	return false
}

/* Obtiene el tipo de elemento */
func (this DiagramBpmnIO) GetTypeElement(elem *etree.Element) (string) {
	if this.isGateway(elem) {
		return BPMNIO_TYPE_GATEWAY
	}
	if this.isEvent(elem) {
		return BPMNIO_TYPE_EVENT
	}
	if this.isTask(elem) {
		return BPMNIO_TYPE_TASK
	}

	return BPMNIO_TYPE_NONE
}

/*
	Esta funcion es la encargada de obtener el elemento principal del diagrama XML
	es decir la etiqueta TAG_ROOT

	Retorna un puntero Element
 */
func (this DiagramBpmnIO) getRootElement() (*etree.Element) {
	return this.documentXML.SelectElement(BPMNIO_TAG_ROOT)
}

/*
	Esta funcion es la encargada de obtener los elementos process donde estan ubicados todas
	las actividades del diagrama, es hija del elemento TAG_ROOT > TAG_PROCESS

	Retorna un slice puntero Element
 */
func (this DiagramBpmnIO) getProcessElements() ([]*etree.Element) {
	return this.getRootElement().SelectElements(BPMNIO_TAG_PROCESS)
}

/*
	Esta funcion es usada para buscar cualquier elemento dentro de la etiquta TAG_PROCESS
	que coincida con el parametro id

	Retorna un puntero Element
 */
func (this DiagramBpmnIO) getElementByID(id string) (*etree.Element) {
	return this.getRootElement().FindElement(`//[@id='` + id + `']`)
}

/*
	Esta funcion es usada para buscar cualquier elemento dentro de la etiquta TAG_ROOT
	que coincida con el atributo proporcionado

	Retorna un puntero Element
 */
func (this DiagramBpmnIO) getElementByAttr(atrib string, val string) (*etree.Element) {
	return this.getRootElement().FindElement(`//[@` + atrib + `='` + val + `']`)
}

/*
	Esta funcion es usada para verificar si un elemento posee datos de entrada

	Retorna booleano
 */
func (this DiagramBpmnIO) HasDataInput(elem *etree.Element) (bool) {
	return len(elem.SelectElements(BPMNIO_TAG_DATA_INPUT_ASSOCIATION)) > 0
}

/*
	Esta funcion obtiene todos los datos de entrada de un elemento

	Retorna slice de puntero element
 */
func (this DiagramBpmnIO) GetDataInputElement(elem *etree.Element) ([]*etree.Element) {
	data := elem.SelectElements(BPMNIO_TAG_DATA_INPUT_ASSOCIATION)
	slice_data := make([]*etree.Element, 0)

	for _, input_asoc := range data {
		ref := input_asoc.SelectElement(BPMNIO_ATTR_SOURCE_REF)
		if ref != nil {
			id_object_ref := ref.Text()
			slice_data = append(slice_data, this.getElementByID(id_object_ref))
		}
	}
	return slice_data
}

/*
	Esta funcion obtiene todos dentro de un lane especifico

	Retorna slice de puntero element
 */
func (this DiagramBpmnIO) GetElementsInLane(lane_elem *etree.Element) ([]*etree.Element) {
	elems_lane := make([]*etree.Element, 0)
	nodes_ref := lane_elem.SelectElements(BPMNIO_TAG_FLOW_NODE_REF)

	for _, node := range nodes_ref {
		elem := this.getElementByID(node.Text())

		if elem != nil {
			elems_lane = append(elems_lane, elem)
		}
	}
	return elems_lane
}

/*
	Esta funcion obtiene el atributo nombre del elemento

	Retorna un string con el nombre
 */
func (this DiagramBpmnIO) GetAttributeElement(elem *etree.Element, key string) (string) {
	return elem.SelectAttrValue(key, EMPTY)
}

/*
	Esta funcion obtiene todoo el flujo del proceso, secuencial en un slice de elementos

	Retorna slice de puntero element
 */
func (this DiagramBpmnIO) GetSuccessionProcess() []*etree.Element {
	//ojo Sequence Element
	var sq_elem, elem_add *etree.Element
	s := make([]*etree.Element, 0)

	sq_elem = this.getBeginElement()

	elem_add = this.getElementByID(sq_elem.SelectAttrValue(BPMNIO_ATTR_SOURCE_REF, ""))
	s = append(s, elem_add)

	for {
		if this.hasMoreElements(sq_elem) {
			sq_elem = this.getNextElement(sq_elem)

			elem_add = this.getElementByID(sq_elem.SelectAttrValue(BPMNIO_ATTR_SOURCE_REF, ""))
			s = append(s, elem_add)

		} else {
			elem_add = this.getElementByID(sq_elem.SelectAttrValue(BPMNIO_ATTR_TARGET_REF, ""))
			s = append(s, elem_add)
			break
		}
	}
	return s
}

/*
	Esta funcion obtiene todos los carriles del diagrama

	Retorna slice de puntero element
 */
func (this DiagramBpmnIO) GetLanesElement() []*etree.Element {
	lanes := make([]*etree.Element, 0)

	for _, e_process := range this.getProcessElements() {
		lane_set := e_process.SelectElement(BPMNIO_TAG_LANE_SET)

		if lane_set != nil {
			for _, lane := range lane_set.SelectElements(BPMNIO_TAG_LANE) {
				lanes = append(lanes, lane)
			}
		}
	}

	return lanes
}

/*
	Esta funcion obtiene el elemento principal basandose en TAG_START_EVENT

	Retorna un puntero element
 */
func (this DiagramBpmnIO) getBeginElement() (*etree.Element) {

	for _, process := range this.getProcessElements() {
		start_event := process.SelectElement(BPMNIO_TAG_START_EVENT)

		if start_event != nil {
			for _, flow := range this.flows {
				if flow.SelectAttr(BPMNIO_ATTR_SOURCE_REF).Value == start_event.SelectAttr(BPMNIO_ATTR_ID).Value {
					return flow
				}
			}
		}
	}

	return nil
}

/*
	Esta funcion obtiene el siguiente elemento en la sequencia del proceso

	Retorna un puntero element o nil
 */
func (this DiagramBpmnIO) getNextElement(flow_previus *etree.Element) (*etree.Element) {
	for _, seq := range this.flows {
		if seq.SelectAttr(BPMNIO_ATTR_SOURCE_REF).Value == flow_previus.SelectAttr(BPMNIO_ATTR_TARGET_REF).Value {
			return seq
		}
	}
	return nil
}

/*
	Esta funcion verifica si existen mas elementos en la sequencia del proceso

	Retorna booleano
 */
func (this DiagramBpmnIO) hasMoreElements(flow_previus *etree.Element) bool {
	return this.getNextElement(flow_previus) != nil
}

/*
	Esta funcion obtiene los flows del proceso

	Retorna Retorna slice de puntero element
 */
func (this DiagramBpmnIO) GetFlows() ([]*etree.Element) {
	return this.flows
}

/*
	Esta funcion obtiene todos los ids de los elementos del diagrama..

	Retorna Retorna slice de strings
 */
func (this DiagramBpmnIO) getIDsElements() ([]string) {
	elements_ids := make([]string, 0)

	for _, flow := range this.GetFlows() {
		id_element_source := flow.SelectAttrValue(BPMNIO_ATTR_SOURCE_REF, EMPTY)
		id_element_target := flow.SelectAttrValue(BPMNIO_ATTR_TARGET_REF, EMPTY)

		if !utils.InSlice(id_element_source, elements_ids) {
			elements_ids = append(elements_ids, id_element_source)
		}

		if !utils.InSlice(id_element_target, elements_ids) {
			elements_ids = append(elements_ids, id_element_target)
		}
	}

	return elements_ids
}

/*
	Esta funcion obtiene todos los elementos del diagrama ya sean gateways, process, subprocess etc..

	Retorna Retorna slice de puntero element
 */
func (this DiagramBpmnIO) GetElements() ([]*etree.Element) {
	elements := make([]*etree.Element, 0)

	for _, id := range this.getIDsElements() {
		elements = append(elements, this.getElementByID(id))
	}

	return elements
}

/***********************************************/
//   functions for interface Diagram
/***********************************************/

func (this DiagramBpmnIO) GetGateways() ([]Gateway) {
	s_gateways := make([]Gateway, 0)
	return s_gateways
}

func (this DiagramBpmnIO) GetEvents() ([]Event) {
	events := make([]Event, 0)
	for _, elem := range this.GetElements() {
		if this.isEvent(elem) {
			events = append(events, Event{
				Type: this.GetTypeElement(elem),
			})
		}
	}
	return events
}

func (this DiagramBpmnIO) GetTasks() ([]Task) {
	tasks := make([]Task, 0)
	for _, elem := range this.GetElements() {
		if this.isTask(elem) {
			tasks = append(tasks, Task{
				Name:     elem.SelectAttrValue(BPMNIO_ATTR_NAME, EMPTY),
				Type:     this.GetTypeElement(elem),
				NeuronID: this.GetAttributeElement(elem, BPMNIO_ATTR_CVDI_NEURON),
				ActionID: this.GetAttributeElement(elem, BPMNIO_ATTR_CVDI_ACTION),
			})
		}
	}
	return tasks
}

func (this DiagramBpmnIO) GetLanes() ([]Lane) {
	lanes := make([]Lane, 0)
	for _, lane := range this.GetLanesElement() {
		lanes = append(lanes, Lane{
			Name: lane.SelectAttrValue(BPMNIO_ATTR_NAME, EMPTY),
		})
	}
	return lanes
}

func (this DiagramBpmnIO) GetPools() ([]Pool) {
	return make([]Pool, 0)
}

func (this *DiagramBpmnIO) LoadDiagramByPath(path string) error {
	return this.ReadFromFile(path)
}

func (this *DiagramBpmnIO) LoadDiagramByBuffer(buf []byte) error {
	return this.ReadFromBytes(buf)
}

func (this *DiagramBpmnIO) LoadDiagramByString(str string) error {
	return this.ReadFromString(str)
}
